package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/webitel/logger/model"

	errors "github.com/webitel/engine/model"
)

type Storage interface {
	// Interface to the log table
	Log() LogStore
	// Interface to the config table
	Config() ConfigStore
	// Database connection
	Database() (*sqlx.DB, errors.AppError)
	// Initializes logger schema
	//SchemaInit() errors.AppError
	// Opens connection to the storage
	Open() errors.AppError
	// Closes connection to the storage
	Close() errors.AppError
}

//var RecordTablesMapping = map[string]TableProperties{
//	"schema":   {Schema: "flow", Table: "acr_routing_scheme", ColumnDomain: "domain_id", ColumnName: "name"},
//	"cc_queue": {Schema: "call_center", Table: "cc_queue", ColumnDomain: "domain_id", ColumnName: "name"},
//}
//
//type TableProperties struct {
//	Schema       string
//	Table        string
//	ColumnDomain string
//	ColumnName   string
//}
//
//func (t *TableProperties) GetFullPath() string {
//	return fmt.Sprintf("%s.%s", t.Schema, t.Table)
//}

type LogStore interface {
	Insert(ctx context.Context, log *model.Log, domainId int) errors.AppError
	//GetByObjectId(ctx context.Context, opt *model.SearchOptions, domainId int, objectId int) (*[]model.Log, errors.AppError)
	//GetByObjectIdWithDates(ctx context.Context, domainId int, objectId int, dateFrom time.Time, dateTo time.Time) (*[]model.Log, errors.AppError)
	//GetByConfigId(ctx context.Context, opt *model.SearchOptions, configId int) (*[]model.Log, errors.AppError)
	//GetByConfigIdWithDates(ctx context.Context, configId int, dateFrom time.Time, dateTo time.Time) (*[]model.Log, errors.AppError)
	//GetByUserId(ctx context.Context, opt *model.SearchOptions, userId int) (*[]model.Log, errors.AppError)
	Get(ctx context.Context, opt *model.SearchOptions, filters any) ([]*model.Log, errors.AppError)
	InsertMany(ctx context.Context, log []*model.Log, domainId int) errors.AppError
	DeleteByLowerThanDate(ctx context.Context, date time.Time, configId int) (int, errors.AppError)
}

type ConfigStore interface {
	CheckAccess(ctx context.Context, domainId, id int64, groups []int, access uint32) (bool, errors.AppError)
	// GetAvailableSystemObjects - get all available objects from domain which are named as [filters]
	GetAvailableSystemObjects(ctx context.Context, domainId int, includeExisting bool, filters ...string) ([]*model.Lookup, errors.AppError)
	//CheckAccessByObjectId(ctx context.Context, domainId, objectId int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError)
	Update(ctx context.Context, conf *model.Config, fields []string, userId int) (*model.Config, errors.AppError)
	Insert(ctx context.Context, conf *model.Config, userId int) (*model.Config, errors.AppError)
	Get(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, filters any) ([]*model.Config, errors.AppError)
	GetByObjectId(ctx context.Context, domainId int, objectId int) (*model.Config, errors.AppError)
	//GetAll(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, domainId int) (*[]model.Config, errors.AppError)
	//GetAllEnabledConfigs(ctx context.Context) (*[]model.Config, errors.AppError)
	GetById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, errors.AppError)
	Delete(ctx context.Context, id int32) errors.AppError
	DeleteMany(ctx context.Context, rbac *model.RbacOptions, ids []int32) errors.AppError
}
