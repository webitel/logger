package storage

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"

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
	CheckRecordExist(ctx context.Context, objectName string, recordId int32) (bool, errors.AppError)
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
type Table struct {
	Path       string
	NameColumn string
}

// ApplyFiltersToBuilder determines type of filters parameter and applies filters to the base according to the determined type.
// columnAlias is additional parameter applied to every model.Filter existing in filters and checks if model.Filter.Column has alias in the {columnAlias}
func ApplyFiltersToBuilderBulk(base any, columnAlias map[string]string, filters any) (any, errors.AppError) {
	if filters == nil {
		return base, nil
	}
	switch data := filters.(type) {
	case *model.FilterNode:
		switch data.Connection {
		case model.AND:
			result := squirrel.And{}
			for _, bunch := range data.Nodes {
				switch bunchType := bunch.(type) {
				case *model.FilterNode:
					lowerResult, err := ApplyFiltersToBuilderBulk(result, columnAlias, bunchType)
					if err != nil {
						return nil, err
					}
					switch newData := lowerResult.(type) {
					case squirrel.And:
						result = append(result, newData)
					}
				case *model.Filter:
					result = append(result, applyFilter(bunchType, columnAlias))
				}
			}

			switch baseType := base.(type) {
			case squirrel.And:
				base = append(baseType, result)
			case squirrel.Or:
				base = append(baseType, result)
			case squirrel.SelectBuilder:
				base = baseType.Where(result)
			}
			return base, nil
		case model.OR:
			result := squirrel.Or{}
			for _, bunch := range data.Nodes {
				switch v := bunch.(type) {
				case *model.FilterNode:
					lowerResult, err := ApplyFiltersToBuilderBulk(result, columnAlias, v)
					if err != nil {
						return nil, err
					}
					switch newData := lowerResult.(type) {
					case squirrel.And:
						result = append(result, newData)
					}
				case *model.Filter:
					result = append(result, applyFilter(v, columnAlias))
				}
			}
			switch baseType := base.(type) {
			case squirrel.And:
				base = append(baseType, result)
			case squirrel.Or:
				base = append(baseType, result)
			case squirrel.SelectBuilder:
				base = baseType.Where(result)
			}
			return base, nil
		}
	case *model.Filter:
		switch baseType := base.(type) {
		case squirrel.And:
			base = append(baseType, applyFilter(data, columnAlias))
		case squirrel.Or:
			base = append(baseType, applyFilter(data, columnAlias))
		case squirrel.SelectBuilder:
			base = baseType.Where(applyFilter(data, columnAlias))
		}
	}

	return base, nil
}

// Apply filter performs convertation between model.Filter and squirrel.Sqlizer.
// columnAlias is additional parameter to determine if model.Filter in the Column property has alias of the column and NOT the real DB column name.
func applyFilter(filter *model.Filter, columnsAlias map[string]string) squirrel.Sqlizer {
	columnName := filter.Column
	if columnsAlias != nil {
		if alias, ok := columnsAlias[columnName]; ok {
			columnName = alias
		}
	}
	var result squirrel.Sqlizer
	switch filter.ComparisonType {
	case model.GreaterThan:
		result = squirrel.Gt{columnName: filter.Value}
	case model.GreaterThanOrEqual:
		result = squirrel.GtOrEq{columnName: filter.Value}
	case model.LessThan:
		result = squirrel.Lt{columnName: filter.Value}
	case model.LessThanOrEqual:
		result = squirrel.LtOrEq{columnName: filter.Value}
	case model.NotEqual:
		result = squirrel.NotEq{columnName: filter.Value}
	case model.Like:
		result = squirrel.Like{columnName: filter.Value}
	case model.ILike:
		result = squirrel.ILike{columnName: filter.Value}
	default:
		result = squirrel.Eq{columnName: filter.Value}
	}
	return result
}
