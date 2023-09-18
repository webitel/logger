package watcher

import (
	"fmt"
	"time"

	"github.com/webitel/wlog"
)

type WatcherNotify func()

type Watcher struct {
	name            string
	stop            chan struct{}
	stopped         chan struct{}
	pollingInterval int64
	PollAndNotify   WatcherNotify
}

func MakeWatcher(name string, pollingInterval int64, pollAndNotify WatcherNotify) *Watcher {
	return &Watcher{
		name:            name,
		stop:            make(chan struct{}),
		stopped:         make(chan struct{}),
		pollingInterval: pollingInterval,
		PollAndNotify:   pollAndNotify,
	}
}

func (watcher *Watcher) Start() {
	wlog.Debug(fmt.Sprintf("watcher [%s] started", watcher.name))
	//<-time.After(time.Duration(rand.Intn(watcher.pollingInterval)) * time.Millisecond)

	defer func() {
		wlog.Debug(fmt.Sprintf("watcher [%s] finished", watcher.name))
		close(watcher.stopped)
	}()

	for {
		wlog.Info(fmt.Sprintf("watcher [%s] will run after %d seconds", watcher.name, watcher.pollingInterval))
		select {
		case <-watcher.stop:
			wlog.Info(fmt.Sprintf("watcher [%s] received stop signal", watcher.name))
			return
		case <-time.After(time.Duration(watcher.pollingInterval)):
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
