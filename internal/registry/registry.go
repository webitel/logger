package registry

import (
	"github.com/webitel/logger/internal/model"
	"time"
)

const (
	DeregisterCriticalServiceAfter = 30 * time.Second
	ServiceName                    = "logger"
	CheckInterval                  = 5 * time.Second
)

type ServiceRegistrar interface {
	Register() model.AppError
	Deregister() model.AppError
}
