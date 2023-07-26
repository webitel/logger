package app

import (
	"context"
	"fmt"
	errors "github.com/webitel/engine/model"
	"time"
	"webitel_logger/watcher"
)

func (a *App) initializeWatchers() errors.AppError {
	err := a.initializeLogCleaners()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initializeLogCleaners() errors.AppError {
	a.watchers = make(map[string]*watcher.Watcher)
	configs, appErr := a.storage.Config().GetAllEnabledConfigs(context.Background())
	if appErr != nil {
		return appErr
	}
	for _, v := range *configs {
		name := FormatConfigKey(DeleteWatcherPrefix, v.DomainId, v.ObjectId)
		a.GetWatcherByKey(name)
		a.watchers[name] = watcher.MakeWatcher(name, int(time.Hour)*24, a.BuildWatcherDeleteFunction(v.DomainId, v.ObjectId, v.DaysToStore))
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
func (a *App) UpdateDeleteWatcherWithNewInterval(domainId, objectId, dayseToStore int) {
	name := FormatConfigKey(DeleteWatcherPrefix, domainId, objectId)
	val, ok := a.watchers[name]
	if !ok {
		return
	}
	val.Stop()
	val.PollAndNotify = a.BuildWatcherDeleteFunction(domainId, objectId, dayseToStore)
	go val.Start()
}

// New interval in days
func (a *App) InsertNewDeleteWatcher(domainId, objectId, dayseToStore int) {
	name := FormatConfigKey(DeleteWatcherPrefix, domainId, objectId)
	a.watchers[name] = watcher.MakeWatcher(name, int(time.Hour)*24, a.BuildWatcherDeleteFunction(domainId, objectId, dayseToStore))
}

func FormatConfigKey(prefix string, domainId int, objectId int) string {
	return fmt.Sprintf("%s.%d.%d", prefix, domainId, objectId)
}

func (a *App) BuildWatcherDeleteFunction(domainId, objectId, daysToStore int) watcher.WatcherNotify {
	name := FormatConfigKey(DeleteWatcherPrefix, domainId, objectId)
	return func() {
		res, err := a.storage.Log().DeleteByLowerThanDate(context.Background(), time.Now().AddDate(0, 0, -daysToStore), domainId, objectId)
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
