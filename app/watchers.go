package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	proto "github.com/webitel/protos/logger"

	"github.com/webitel/wlog"

	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/watcher"
)

// region COMMON
func (a *App) initializeWatchers() errors.AppError {

	configs, appErr := a.storage.Config().Get(
		context.Background(),
		nil,
		nil,
		model.Filter{
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
	a.uploadWatchers, a.deleteWatchers = make(map[string]*watcher.UploadWatcher), make(map[string]*watcher.Watcher)
	for _, config := range configs {
		a.InsertLogCleaner(config.Id, nil, config.DaysToStore)
		params := &watcher.UploadWatcherParams{
			StorageId:    config.Storage.Id.Int(),
			Period:       config.Period,
			NextUploadOn: config.NextUploadOn.Time(),
			LastLogId:    config.LastUploadedLog.Int(),
			UserId:       config.UpdatedBy.Int(),
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
func (a *App) InsertLogCleaner(configId int, startParams *watcher.StartParams, dayseToStore int) {
	name := FormatKey(DeleteWatcherPrefix, configId)
	a.deleteWatchers[name] = watcher.MakeWatcher(name, startParams, &watcher.CustomExecutionParams{ExecuteImmediately: true}, time.Hour*24, a.BuildLogCleanerFunction(configId, dayseToStore))
	go a.deleteWatchers[name].Start()
}

func (a *App) GetLogCleaner(configId int) *watcher.Watcher {
	key := FormatKey(DeleteWatcherPrefix, configId)
	val, ok := a.deleteWatchers[key]
	if !ok {
		return nil
	}
	return val
}

func (a *App) DeleteLogCleaner(configId ...int) {
	for _, s := range configId {
		key := FormatKey(DeleteWatcherPrefix, s)
		val, ok := a.deleteWatchers[key]
		if !ok {
			return
		}
		val.Stop()
		delete(a.deleteWatchers, key)
	}

}

func (a *App) UpdateLogCleanerWithNewInterval(configId, dayseToStore int) {
	name := FormatKey(DeleteWatcherPrefix, configId)
	val, ok := a.deleteWatchers[name]
	if !ok {
		return
	}
	val.Stop()
	delete(a.deleteWatchers, name)
	a.InsertLogCleaner(configId, nil, dayseToStore)
}

func (a *App) BuildLogCleanerFunction(configId, daysToStore int) watcher.WatcherNotify {
	name := FormatKey(DeleteWatcherPrefix, configId)
	return func() {
		res, err := a.storage.Log().DeleteByLowerThanDate(context.Background(), time.Now().AddDate(0, 0, -daysToStore), configId)
		if err != nil {
			wlog.Info(fmt.Sprintf("watcher [%s]: %s", name, err.Error()))
		} else {
			wlog.Info(fmt.Sprintf("watcher [%s]: cleaned %d rows", name, res))
		}
	}
}

// endregion

// region LOG UPLOADER

func (a *App) DeleteLogUploader(configId ...int) {
	for _, s := range configId {
		key := FormatKey(UploadWatcherPrefix, s)
		val, ok := a.uploadWatchers[key]
		if !ok {
			return
		}
		val.Stop()
		delete(a.uploadWatchers, key)
	}

}

func (a *App) GetLogUploader(configId int) *watcher.UploadWatcher {
	key := FormatKey(UploadWatcherPrefix, configId)
	val, ok := a.uploadWatchers[key]
	if !ok {
		return nil
	}
	return val
}

func (a *App) GetLogUploaderParams(configId int) *watcher.UploadWatcherParams {
	name := FormatKey(UploadWatcherPrefix, configId)
	val, ok := a.uploadWatchers[name]
	if !ok {
		return nil
	}
	return val.Params
}

func (a *App) InsertLogUploader(configId int, startParams *watcher.StartParams, params *watcher.UploadWatcherParams) {
	name := FormatKey(UploadWatcherPrefix, configId)
	a.uploadWatchers[name] = watcher.MakeUploadWatcher(name, startParams, &watcher.CustomExecutionParams{ExecuteImmediately: true}, time.Hour*24, params, a.BuildWatcherUploadFunction(configId, params))
	go a.uploadWatchers[name].Start()
}

func (a *App) BuildWatcherUploadFunction(configId int, params *watcher.UploadWatcherParams) watcher.WatcherNotify {
	format := func(text string) string {
		return fmt.Sprintf("watcher - %s: %s", FormatKey(UploadWatcherPrefix, configId), text)
	}
	return func() {
		if time.Now().UTC().Unix() >= params.NextUploadOn.UTC().Unix() {
			filters := []*model.Filter{
				{
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
			logs, appErr := a.storage.Log().Get(context.Background(), &model.SearchOptions{
				Sort: "-id",
			}, model.FilterBunch{
				Bunch:          filters,
				ConnectionType: model.AND,
			})
			if appErr != nil {
				if !IsErrNoRows(appErr) {
					wlog.Info(format(appErr.Error()))
					return
				}
				wlog.Info(format("no new logs..."))
				return
			}
			convertedLogs, appErr := convertLogModelToMessageBulk(logs)
			if appErr != nil {
				wlog.Info(format(appErr.Error()))
				return
			}
			buf := &bytes.Buffer{}
			s := proto.Logs{Items: convertedLogs}
			encodeResult, err := json.Marshal(s)
			if err != nil {
				wlog.Info(format(err.Error()))
				return
			}
			_, err = buf.Write(encodeResult)
			if err != nil {
				wlog.Info(format(err.Error()))
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
				wlog.Info(format(err.Error()))
				return
			}
			lastLogId := logs[0].Id
			nextUpload := calculateNextPeriodFromNow(int32(params.Period))

			params.NextUploadOn = calculateNextPeriodFromNow(int32(params.Period))
			params.LastLogId = lastLogId

			nullLogId, err := model.NewNullInt(convertedLogs[0].Id)
			if err != nil {
				wlog.Info(format(err.Error()))
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
				wlog.Info(format(appErr.Error()))
				return
			}
			wlog.Info(format("logs successfully uploaded to the storage!"))
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

// endregion
