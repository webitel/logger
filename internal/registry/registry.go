package registry

import (
	"time"
)

const (
	DeregisterCriticalServiceAfter = 30 * time.Second
	ServiceName                    = "logger"
	CheckInterval                  = 5 * time.Second
)

type ServiceRegistrar interface {
	Register() error
	Deregister() error
}
