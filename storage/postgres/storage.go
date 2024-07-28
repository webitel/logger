package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/storage"
	"github.com/webitel/wlog"

	_ "github.com/jackc/pgx/stdlib"
)

type PostgresStore struct {
	config            *model.DatabaseConfig
	conn              *sqlx.DB
	logStore          storage.LogStore
	configStore       storage.ConfigStore
	loginAttemptStore storage.LoginAttemptStore
}

func New(config *model.DatabaseConfig) *PostgresStore {

	return &PostgresStore{config: config}
}

func (s *PostgresStore) Log() storage.LogStore {
	if s.logStore == nil {
		log, err := newLogStore(s)
		if err != nil {
			return nil
		}
		s.logStore = log
	}
	return s.logStore
}
func (s *PostgresStore) Config() storage.ConfigStore {
	if s.configStore == nil {
		conf, err := newConfigStore(s)
		if err != nil {
			return nil
		}
		s.configStore = conf
	}
	return s.configStore
}
func (s *PostgresStore) LoginAttempt() storage.LoginAttemptStore {
	if s.loginAttemptStore == nil {
		conf, err := newLoginAttemptStore(s)
		if err != nil {
			return nil
		}
		s.loginAttemptStore = conf
	}
	return s.loginAttemptStore
}

func (s *PostgresStore) Database() (*sqlx.DB, model.AppError) {
	if s.conn == nil {
		model.NewInternalError("postgres.storage.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *PostgresStore) Open() model.AppError {
	db, err := sqlx.Connect("pgx", s.config.Url)
	//db, err := sql.Open("pgx", s.config.Url)
	if err != nil {
		return model.NewInternalError("postgres.storage.open.connect.fail", err.Error())
	}
	s.conn = db
	wlog.Debug(fmt.Sprintf("postgres: connection opened"))
	return nil
}

func (s *PostgresStore) Close() model.AppError {
	err := s.conn.Close()
	if err != nil {
		return model.NewInternalError("postgres.storage.close.disconnect.fail", fmt.Sprintf("postgres: %s", err.Error()))
	}
	s.conn = nil
	wlog.Debug(fmt.Sprintf("postgres: connection closed"))
	return nil

}

// ApplyFiltersToBuilder determines type of {filters} parameter and applies {filters} to the {base} according to the determined type.
// columnAlias is additional parameter applied to every model.Filter existing in {filters} and checks if {model.Filter.Column} has alias in the {columnAlias}
//func ApplyFiltersToBuilder(base squirrel.SelectBuilder, columnAlias map[string]string, filters any) squirrel.SelectBuilder {
//	switch data := filters.(type) {
//	case model.FilterArray:
//		switch data.Connection {
//		case model.AND:
//			result := squirrel.And{}
//			for _, bunch := range data.Filters {
//				switch bunch.ConnectionType {
//				case model.AND:
//					lowerResult := squirrel.And{}
//					for _, filter := range bunch.Bunch {
//						lowerResult = append(lowerResult, applyFilter(filter, columnAlias))
//					}
//					result = append(result, lowerResult)
//
//				case model.OR:
//					lowerResult := squirrel.Or{}
//					for _, filter := range bunch.Bunch {
//						lowerResult = append(lowerResult, applyFilter(filter, columnAlias))
//					}
//					result = append(result, lowerResult)
//
//				}
//
//			}
//			base = base.Where(result)
//			return base
//		case model.OR:
//			result := squirrel.Or{}
//			for _, bunch := range data.Filters {
//				switch bunch.ConnectionType {
//				case model.AND:
//					lowerResult := squirrel.And{}
//					for _, filter := range bunch.Bunch {
//						lowerResult = append(lowerResult, applyFilter(filter, columnAlias))
//					}
//					result = append(result, lowerResult)
//					base = base.Where(result)
//					return base
//				case model.OR:
//					lowerResult := squirrel.Or{}
//					for _, filter := range bunch.Bunch {
//						lowerResult = append(lowerResult, applyFilter(filter, columnAlias))
//					}
//					result = append(result, lowerResult)
//					base = base.Where(result)
//					return base
//				}
//
//			}
//			base = base.Where(result)
//			return base
//		}
//	case model.FilterBunch:
//		switch data.ConnectionType {
//		case model.AND:
//			result := squirrel.And{}
//			for _, filter := range data.Bunch {
//				result = append(result, applyFilter(filter, columnAlias))
//			}
//
//			base = base.Where(result)
//			return base
//		case model.OR:
//			result := squirrel.Or{}
//			for _, filter := range data.Bunch {
//				result = append(result, applyFilter(filter, columnAlias))
//			}
//			base = base.Where(result)
//			return base
//		}
//	case model.Filter:
//		base = base.Where(applyFilter(&data, columnAlias))
//	}
//
//	return base
//}

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
