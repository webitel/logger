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
	//SchemaInit() errors.AppError
	// Opens connection to the storage
	Open() errors.AppError
	// Closes connection to the storage
	Close() errors.AppError
}

type LogStore interface {
	Insert(context.Context, *model.Log) (*model.Log, errors.AppError)
	//GetByObjectId(ctx context.Context, opt *model.SearchOptions, domainId int, objectId int) (*[]model.Log, errors.AppError)
	//GetByObjectIdWithDates(ctx context.Context, domainId int, objectId int, dateFrom time.Time, dateTo time.Time) (*[]model.Log, errors.AppError)
	//GetByConfigId(ctx context.Context, opt *model.SearchOptions, configId int) (*[]model.Log, errors.AppError)
	//GetByConfigIdWithDates(ctx context.Context, configId int, dateFrom time.Time, dateTo time.Time) (*[]model.Log, errors.AppError)
	//GetByUserId(ctx context.Context, opt *model.SearchOptions, userId int) (*[]model.Log, errors.AppError)
	Get(ctx context.Context, opt *model.SearchOptions, filters ...model.Filter) (*[]model.Log, errors.AppError)
	DeleteByLowerThanDate(ctx context.Context, date time.Time, configId int) (int, errors.AppError)
}

type ConfigStore interface {
	//CheckAccess(ctx context.Context, domainId, id int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError)
	//CheckAccessByObjectId(ctx context.Context, domainId, objectId int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError)
	Update(ctx context.Context, conf *model.Config, userId int) (*model.Config, errors.AppError)
	Insert(ctx context.Context, conf *model.Config, userId int) (*model.Config, errors.AppError)
	Get(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, filters ...model.Filter) (*[]model.Config, errors.AppError)
	GetByObjectId(ctx context.Context, domainId int, objectId int) (*model.Config, errors.AppError)
	//GetAll(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, domainId int) (*[]model.Config, errors.AppError)
	//GetAllEnabledConfigs(ctx context.Context) (*[]model.Config, errors.AppError)
	GetById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, errors.AppError)
}
