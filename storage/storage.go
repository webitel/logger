package storage

import (
	"context"
	"database/sql"
	"webitel_logger/model"

	errors "github.com/webitel/engine/model"
)

type Storage interface {
	// Interface to the log table
	Log() LogStore
	// Interface to the config table
	Config() ConfigStore
	// Database connection
	Database() (*sql.DB, errors.AppError)
	// Initializes logger schema
	SchemaInit() errors.AppError
	// Opens connection to the storage
	Open() errors.AppError
	// Closes connection to the storage
	Close() errors.AppError
}

type LogStore interface {
	Insert(context.Context, *model.Log) (*model.Log, errors.AppError)
	GetByObjectId(context.Context, int) (*[]model.Log, errors.AppError)
	GetByUserId(context.Context, int) (*[]model.Log, errors.AppError)
}

type ConfigStore interface {
	Update(context.Context, *model.Config) (*model.Config, errors.AppError)
	GetByObjectId(context.Context, int) (*model.Config, errors.AppError)
	GetAll(context.Context) (*[]model.Config, errors.AppError)
}
