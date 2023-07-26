package storage

import (
	"context"
	"database/sql"
	"time"
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
	GetByObjectId(ctx context.Context, opt *model.SearchOptions, domainId int, objectId int) (*[]model.Log, errors.AppError)
	GetByConfigId(ctx context.Context, opt *model.SearchOptions, configId int) (*[]model.Log, errors.AppError)
	GetByUserId(ctx context.Context, opt *model.SearchOptions, userId int) (*[]model.Log, errors.AppError)
	DeleteByLowerThanDate(ctx context.Context, date time.Time, domainId int, objectId int) (int, errors.AppError)
}

type ConfigStore interface {
	Update(context.Context, *model.Config) (*model.Config, errors.AppError)
	Insert(context.Context, *model.Config) (*model.Config, errors.AppError)
	GetByObjectId(ctx context.Context /*opt *model.SearchOptions,*/, domainId int, objectId int) (*model.Config, errors.AppError)
	GetAll(ctx context.Context, opt *model.SearchOptions, domainId int) (*[]model.Config, errors.AppError)
	GetAllEnabledConfigs(ctx context.Context) (*[]model.Config, errors.AppError)
	GetById(ctx context.Context, opt *model.SearchOptions, id int) (*model.Config, errors.AppError)
}
