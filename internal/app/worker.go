package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/webitel/logger/internal/watcher"
	"go.opentelemetry.io/otel/attribute"
	"time"

	"log/slog"

	"github.com/webitel/logger/internal/model"
)

// region COMMON
func (a *App) initializeWatchers() model.AppError {
	configs, appErr := a.storage.Config().Select(
		context.Background(),
		nil,
		nil,
		&model.Filter{
			Column:         "object_config.enabled",
			Value:          true,
			ComparisonType: model.Equal,
		},
	)
	if appErr != nil {
		if !IsErrNoRows(appErr) {
			return appErr
		}
	}
	a.logUploaders, a.logCleaners = make(map[string]*watcher.UploadWatcher), make(map[string]*watcher.Watcher)
	for _, config := range configs {
		a.InsertLogCleaner(config.Id, nil, config.DaysToStore)
		params := &watcher.UploadWatcherParams{
			StorageId:    config.Storage.Id.Int(),
			Period:       config.Period,
			NextUploadOn: config.NextUploadOn.Time(),
			LastLogId:    config.LastUploadedLog.Int(),
			UserId:       config.UpdatedBy.Int64(),
			DomainId:     config.DomainId,
		}
		a.InsertLogUploader(config.Id, nil, params)
	}
	return nil
}

func (a *App) DeleteWatchers(configId ...int) {
	a.DeleteLogUploader(configId...)
	a.DeleteLogCleaner(configId...)
}

// endregion

// region LOG CLEANER

func (a *App) InsertLogCleaner(configId int, startParams *watcher.StarterParams, dayseToStore int) {
	name := FormatKey(DeleteWatcherPrefix, configId)
	a.logCleaners[name] = watcher.NewWatcher(name, startParams, &watcher.CustomExecutionParams{ExecuteImmediately: true}, time.Hour*24, a.BuildLogCleanerFunction(configId, dayseToStore))
	go a.logCleaners[name].Start()
}

func (a *App) GetLogCleaner(configId int) *watcher.Watcher {
	key := FormatKey(DeleteWatcherPrefix, configId)
	val, ok := a.logCleaners[key]
	if !ok {
		return nil
	}
	return val
}

func (a *App) DeleteLogCleaner(configId ...int) {
	for _, s := range configId {
		key := FormatKey(DeleteWatcherPrefix, s)
		val, ok := a.logCleaners[key]
		if !ok {
			return
		}
		val.Stop()
		delete(a.logCleaners, key)
	}

}

func (a *App) UpdateLogCleanerWithNewInterval(ctx context.Context, configId, daysToStore int) {
	ctx, span := a.tracer.Start(ctx, "app.UpdateLogCleaner")
	defer span.End()
	span.SetAttributes(attribute.Int("worker.interval", daysToStore))
	name := FormatKey(DeleteWatcherPrefix, configId)
	val, ok := a.logCleaners[name]
	if !ok {
		return
	}
	val.Stop()
	delete(a.logCleaners, name)
	a.InsertLogCleaner(configId, nil, daysToStore)
	slog.InfoContext(ctx, fmt.Sprintf("[%s]: recreated with new parameters", name))
}

func (a *App) BuildLogCleanerFunction(configId, daysToStore int) watcher.WatcherRoutine {
	name := FormatKey(DeleteWatcherPrefix, configId)
	logAttr := slog.Group(
		"worker",
		slog.String("name", name),
		slog.Int("config_id", configId),
		slog.Int("period", daysToStore),
	)
	return func() {
		res, err := a.DeleteLogs(context.Background(), configId, time.Now().AddDate(0, 0, -daysToStore))
		if err != nil {
			slog.Warn(fmt.Sprintf("[%s]: %s", name, err.Error()), logAttr)
		} else {
			slog.Debug(fmt.Sprintf("[%s]: cleaned %d rows", name, res), logAttr, slog.Int("cleaned", res))
		}
	}
}

// endregion

// region LOG UPLOADER

func (a *App) DeleteLogUploader(configId ...int) {
	for _, s := range configId {
		key := FormatKey(UploadWatcherPrefix, s)
		val, ok := a.logUploaders[key]
		if !ok {
			return
		}
		val.Stop()
		delete(a.logUploaders, key)
	}

}

func (a *App) GetLogUploader(configId int) *watcher.UploadWatcher {
	key := FormatKey(UploadWatcherPrefix, configId)
	val, ok := a.logUploaders[key]
	if !ok {
		return nil
	}
	return val
}

func (a *App) GetLogUploaderParams(configId int) *watcher.UploadWatcherParams {
	name := FormatKey(UploadWatcherPrefix, configId)
	val, ok := a.logUploaders[name]
	if !ok {
		return nil
	}
	return val.Params
}

func (a *App) InsertLogUploader(configId int, startParams *watcher.StarterParams, params *watcher.UploadWatcherParams) {
	name := FormatKey(UploadWatcherPrefix, configId)
	a.logUploaders[name] = watcher.NewUploadWatcher(name, startParams, &watcher.CustomExecutionParams{ExecuteImmediately: true}, params, time.Hour*24, a.BuildWatcherUploadFunction(configId, params))
	go a.logUploaders[name].Start()
}

func (a *App) BuildWatcherUploadFunction(configId int, params *watcher.UploadWatcherParams) watcher.WatcherRoutine {
	name := FormatKey(UploadWatcherPrefix, configId)
	format := func(text string) string {
		return fmt.Sprintf("[%s]: %s", name, text)
	}
	return func() {
		if time.Now().UTC().Unix() >= params.NextUploadOn.UTC().Unix() {
			logAttr := slog.Group(
				"worker",
				slog.String("name", name),
				slog.Int("config_id", configId),
				slog.Int64("domain", params.DomainId),
				slog.Int("period", params.Period),
				slog.Int("storage", params.StorageId),
			)
			filters := []any{
				&model.Filter{
					Column:         "config_id",
					Value:          configId,
					ComparisonType: model.Equal,
				}}
			if v := params.LastLogId; v != 0 {
				filters = append(filters, &model.Filter{
					Column:         "id",
					Value:          v,
					ComparisonType: model.GreaterThan,
				})
			}
			logs, appErr := a.storage.Log().Select(context.Background(), &model.SearchOptions{
				Sort: "-id",
			}, &model.FilterNode{
				Nodes:      filters,
				Connection: model.AND,
			})
			if appErr != nil {
				if !IsErrNoRows(appErr) {
					slog.Warn(format(appErr.Error()), logAttr)
					return
				}
				slog.Info(format("no new logs..."))
				return
			}
			buf := &bytes.Buffer{}
			encodeResult, err := json.Marshal(logs)
			if err != nil {
				slog.Warn(format(err.Error()), logAttr)
				return
			}
			_, err = buf.Write(encodeResult)
			if err != nil {
				slog.Warn(format(err.Error()), logAttr)
				return
			}
			year, month, day := time.Now().Date()
			fileName := fmt.Sprintf("log_%d_%d_%s_%d.json", configId, year, month, day)
			uuid := fmt.Sprintf("%s-%d", logs[0].Object.Name, configId)
			_, err = a.UploadFile(context.Background(), int64(params.DomainId), uuid, params.StorageId, buf, model.File{
				Name:     fileName,
				MimeType: "application/json",
				ViewName: &fileName,
			})
			if err != nil {
				slog.Warn(format(err.Error()), logAttr)
				return
			}
			lastLogId := logs[0].Id
			nextUpload := calculateNextPeriodFromNow(int32(params.Period))

			params.NextUploadOn = calculateNextPeriodFromNow(int32(params.Period))
			params.LastLogId = lastLogId

			nullLogId, err := model.NewNullInt(logs[0].Id)
			if err != nil {
				slog.Warn(format(err.Error()), logAttr)
				return
			}
			_, appErr = a.storage.Config().Update(
				context.Background(),
				&model.Config{
					Id:              configId,
					NextUploadOn:    *model.NewNullTime(nextUpload),
					LastUploadedLog: *nullLogId,
				},
				[]string{"next_upload_on", "last_uploaded_log_id"},
				params.UserId,
			)
			if appErr != nil {
				slog.Warn(format(appErr.Error()), logAttr)
				return
			}
			slog.Info(format("logs successfully uploaded to the storage!"), logAttr)
		}
	}
}

// endregion

// region UTILS
func FormatKey(prefix string, args ...any) string {
	base := prefix
	for _, v := range args {
		base += fmt.Sprintf(".%d", v)
	}
	return base
}

func calculateNextPeriodFromNow(period int32) time.Time {
	now := time.Now().Add(time.Hour * 24 * time.Duration(period))
	return now
}

// endregion
