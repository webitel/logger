package watcher

import (
	"fmt"
	"time"

	"github.com/webitel/wlog"
)

type WatcherNotify func()

type Watcher struct {
	name                  string
	stop                  chan struct{}
	stopped               chan struct{}
	startParams           *StartParams
	customExecutionParams *CustomExecutionParams
	PollAndNotify         WatcherNotify
	interval              time.Duration
}

type UploadWatcher struct {
	Watcher
	Params *UploadWatcherParams
}

type UploadWatcherParams struct {
	StorageId    int
	Period       int
	NextUploadOn time.Time
	LastLogId    int
	DomainId     int64
	UserId       int
}

type StartParams struct {
	StartPollingAfter time.Duration
}

type CustomExecutionParams struct {
	ExecuteImmediately bool
}

func MakeWatcher(name string, startParams *StartParams, customExecutionParams *CustomExecutionParams, interval time.Duration, pollAndNotify WatcherNotify) *Watcher {
	return &Watcher{
		name:                  name,
		stop:                  make(chan struct{}),
		stopped:               make(chan struct{}),
		interval:              interval,
		startParams:           startParams,
		customExecutionParams: customExecutionParams,
		PollAndNotify:         pollAndNotify,
	}
}

func MakeUploadWatcher(name string, startParams *StartParams, customExecutionParams *CustomExecutionParams, pollingInterval time.Duration, params *UploadWatcherParams, pollAndNotify WatcherNotify) *UploadWatcher {
	return &UploadWatcher{
		Watcher: Watcher{
			name:                  name,
			stop:                  make(chan struct{}),
			stopped:               make(chan struct{}),
			interval:              pollingInterval,
			startParams:           startParams,
			customExecutionParams: customExecutionParams,
			PollAndNotify:         pollAndNotify,
		},
		Params: params,
	}
}

func (watcher *Watcher) Start() {
	wlog.Debug(fmt.Sprintf("watcher [%s] started", watcher.name))
	defer func() {
		wlog.Debug(fmt.Sprintf("watcher [%s] finished", watcher.name))
		close(watcher.stopped)
	}()
	if watcher.startParams != nil {
		time.Sleep(watcher.startParams.StartPollingAfter)
	}
	if watcher.customExecutionParams != nil && watcher.customExecutionParams.ExecuteImmediately {
		watcher.PollAndNotify()
	}

	for {
		wlog.Info(fmt.Sprintf("watcher [%s] will run every %d hours", watcher.name, watcher.interval/time.Hour))
		select {
		case <-watcher.stop:
			wlog.Info(fmt.Sprintf("watcher [%s] received stop signal", watcher.name))
			return
		case <-time.After(watcher.interval):
			watcher.PollAndNotify()
		}
	}

}

func (watcher *Watcher) GetName() string {
	return watcher.name
}

func (watcher *Watcher) Stop() {
	wlog.Debug(fmt.Sprintf("watcher [%s] stopping", watcher.name))
	close(watcher.stop)
	<-watcher.stopped
}
