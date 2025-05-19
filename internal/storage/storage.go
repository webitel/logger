package storage

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	"github.com/webitel/logger/internal/model"
)

type Storage interface {
	// Interface to the log table
	Log() LogStore
	// Interface to the config table
	Config() ConfigStore
	// Interface to the config table
	LoginAttempt() LoginAttemptStore
	// Database connection
	Database() (*sqlx.DB, model.AppError)
	// Opens connection to the storage
	Open() model.AppError
	// Closes connection to the storage
	Close() model.AppError
}

type LogStore interface {
	Insert(ctx context.Context, log *model.Log, domainId int) model.AppError
	Select(ctx context.Context, opt *model.SearchOptions, filters any) ([]*model.Log, model.AppError)
	InsertBulk(ctx context.Context, log []*model.Log, domainId int) model.AppError
	Delete(ctx context.Context, earlierThan time.Time, configId int) (int, model.AppError)

	CheckRecordExist(ctx context.Context, objectName string, recordId int32) (bool, model.AppError)
}

type ConfigStore interface {
	// GetAvailableSystemObjects - get all available objects from domain which are named as [filters]
	Insert(ctx context.Context, conf *model.Config, userId int64) (*model.Config, model.AppError)
	Select(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, filters any) ([]*model.Config, model.AppError)
	Update(ctx context.Context, conf *model.Config, fields []string, userId int64) (*model.Config, model.AppError)
	Delete(ctx context.Context, id int32, domainId int64) model.AppError

	GetAvailableSystemObjects(ctx context.Context, domainId int, includeExisting bool, filters ...string) ([]*model.Lookup, model.AppError)
	CheckAccess(ctx context.Context, domainId, id int64, groups []int64, access uint8) (bool, model.AppError)

	GetByObjectId(ctx context.Context, domainId int, objectId int) (*model.Config, model.AppError)
	GetById(ctx context.Context, rbac *model.RbacOptions, id int, domainId int64) (*model.Config, model.AppError)
	DeleteMany(ctx context.Context, rbac *model.RbacOptions, ids []int32, domainId int64) model.AppError
}

type LoginAttemptStore interface {
	Insert(ctx context.Context, m *model.LoginAttempt) (*model.LoginAttempt, model.AppError)
	Select(ctx context.Context, searchOpts *model.SearchOptions, filters any) ([]*model.LoginAttempt, model.AppError)
}

type Table struct {
	Path       string
	NameColumn string
}

// ApplyFiltersToBuilder determines type of filters parameter and applies filters to the base according to the determined type.
// columnAlias is additional parameter applied to every model.Filter existing in filters and checks if model.Filter.Column has alias in the {columnAlias}
func ApplyFiltersToBuilderBulk(base any, columnAlias map[string]string, filters any) (any, model.AppError) {
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
	default:
		return nil, model.NewInternalError("storage.storage.apply_filters_to_builder_bulk.switch_filter.unknown", "invalid filter type")
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
