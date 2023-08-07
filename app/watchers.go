package app

import (
	"context"
	"fmt"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/watcher"
	"time"
)

func (a *App) initializeWatchers() errors.AppError {
	err := a.initializeLogCleaners()
	if err != nil {
		if !IsErrNoRows(err) {
			return err
		}
	}
	return nil
}

func (a *App) initializeLogCleaners() errors.AppError {
	a.watchers = make(map[string]*watcher.Watcher)
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
		return appErr
	}
	for _, v := range *configs {
		name := FormatKey(DeleteWatcherPrefix, v.DomainId, v.Object.Id.Int())
		a.GetWatcherByKey(name)
		a.watchers[name] = watcher.MakeWatcher(name, int(time.Hour)*24, a.BuildWatcherDeleteFunction(v.Id, v.DaysToStore))
		go a.watchers[name].Start()
	}
	return nil
}

func (a *App) initializeStorageUploaders() errors.AppError {
	//panic("unimplemented")
	//connections, appErr := a.serviceDiscovery.GetByName("storage")
	return nil
}

func (a *App) DeleteWatcherByKey(key string) {
	val, ok := a.watchers[key]
	if !ok {
		return
	}
	val.Stop()
	delete(a.watchers, key)
}

func (a *App) GetWatcherByKey(key string) *watcher.Watcher {
	val, ok := a.watchers[key]
	if !ok {
		return nil
	}
	return val
}

// New interval in days
func (a *App) UpdateDeleteWatcherWithNewInterval(configId, dayseToStore int) {
	name := FormatKey(DeleteWatcherPrefix, configId)
	val, ok := a.watchers[name]
	if !ok {
		return
	}
	val.Stop()
	val.PollAndNotify = a.BuildWatcherDeleteFunction(configId, dayseToStore)
	go val.Start()
}

// New interval in days
func (a *App) InsertNewDeleteWatcher(configId, dayseToStore int) {
	name := FormatKey(DeleteWatcherPrefix, configId)
	a.watchers[name] = watcher.MakeWatcher(name, int(time.Hour)*24, a.BuildWatcherDeleteFunction(configId, dayseToStore))
}

//func FormatCacheKey(prefix string, domainId int, objectId int) string {
//	return fmt.Sprintf("%s.%d.%d", prefix, domainId, objectId)
//}

func FormatKey(prefix string, args ...any) string {
	base := prefix
	for _, v := range args {
		base += fmt.Sprintf(".%s", v)
	}
	return base
}

func (a *App) BuildWatcherDeleteFunction(configId, daysToStore int) watcher.WatcherNotify {
	name := FormatKey(DeleteWatcherPrefix, configId)
	return func() {
		res, err := a.storage.Log().DeleteByLowerThanDate(context.Background(), time.Now().AddDate(0, 0, -daysToStore), configId)
		if err != nil {
			fmt.Printf("error while executing watcher function. watcher - %s, error: %s", name, err.Error())
		}
		fmt.Printf("watcher - %s cleaned %d rows", name, res)
	}
}

func (a *App) BuildWatcherUploadFunction(domainId, objectId, daysToStore int) watcher.WatcherNotify {
	panic("unimplemented")
	//name := FormatConfigKey(DeleteWatcherPrefix, domainId, objectId)
	//return func() {
	//
	//	if err != nil {
	//		wlog.Debug(fmt.Sprintf("error while executing watcher function. watcher - %s, error: %s", name, err.Error()))
	//	}
	//	wlog.Debug(fmt.Sprintf("watcher - %s cleaned %d rows", name, res))
	//}
}
