package watcher

import (
	"fmt"
	"time"

	"log/slog"
)

type WatcherRoutine func()

// Watcher is the representative type for CRON jobs
type Watcher struct {
	// name represents the identifier of the Watcher
	name string
	// stop represents stop-channel for the Watcher, used when Watcher.Stop() called
	// cancels the execution of the routine
	stop chan any
	// startParams represents Watcher starter parameters
	startParams *StarterParams
	// startParams represents Watcher custom parameters
	customExecutionParams *CustomExecutionParams
	// routine represents the function called every interval
	routine WatcherRoutine
	// interval represents time between routine calls
	interval time.Duration
}

// UploadWatcher is the representative type for logs upload CRON job
type UploadWatcher struct {
	Watcher
	Params *UploadWatcherParams
}

// UploadWatcherParams represents attributes required by UploadWatcher
type UploadWatcherParams struct {
	StorageId    *int
	Period       int
	NextUploadOn *time.Time
	LastLogId    *int
	DomainId     int
	UserId       *int
}

// StarterParams represents attributes to customize Watcher start
type StarterParams struct {
	// StartPollingAfter represents the time for watcher to wait given amount of time before start the countdown for the execution
	StartPollingAfter time.Duration
}

// CustomExecutionParams represents attributes to customize Watcher.routine execution
type CustomExecutionParams struct {
	// ExecuteImmediately represents ability to execute routine after watcher started countdown (if StarterParams.StartPollingAfter was set it will wait the duration of StarterParams.StartPollingAfter)
	ExecuteImmediately bool
}

// NewWatcher constructs new watcher
func NewWatcher(name string, startParams *StarterParams, customExecutionParams *CustomExecutionParams, interval time.Duration, routine WatcherRoutine) *Watcher {
	return &Watcher{
		name:                  name,
		stop:                  make(chan any),
		interval:              interval,
		startParams:           startParams,
		customExecutionParams: customExecutionParams,
		routine:               routine,
	}
}

// NewWatcher constructs new logs upload watcher
func NewUploadWatcher(name string, startParams *StarterParams, customExecutionParams *CustomExecutionParams, params *UploadWatcherParams, pollingInterval time.Duration, pollAndNotify WatcherRoutine) *UploadWatcher {
	return &UploadWatcher{
		Watcher: Watcher{
			name:                  name,
			stop:                  make(chan any),
			interval:              pollingInterval,
			startParams:           startParams,
			customExecutionParams: customExecutionParams,
			routine:               pollAndNotify,
		},
		Params: params,
	}
}

func (w *Watcher) Start() {
	logAttr := slog.Group("worker", slog.String("name", w.name), slog.Duration("interval", w.interval))
	slog.Debug(fmt.Sprintf("[%s] started", w.name), logAttr)
	defer func() {
		slog.Debug(fmt.Sprintf("[%s] finished", w.name), logAttr)
		close(w.stop)
	}()
	if w.startParams != nil {
		time.Sleep(w.startParams.StartPollingAfter)
	}
	if w.customExecutionParams != nil && w.customExecutionParams.ExecuteImmediately {
		w.routine()
	}

	for {
		slog.Debug(fmt.Sprintf("[%s] routine call", w.name), logAttr)
		select {
		case <-w.stop:
			return
		case <-time.After(w.interval):
			w.routine()
		}
	}

}

func (w *Watcher) GetName() string {
	return w.name
}

func (w *Watcher) Stop() {
	w.stop <- "close"
}
