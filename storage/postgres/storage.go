package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/storage"

	_ "github.com/jackc/pgx/stdlib"
	errors "github.com/webitel/engine/model"
)

type PostgresStore struct {
	config      *model.DatabaseConfig
	conn        *sqlx.DB
	logStore    storage.LogStore
	configStore storage.ConfigStore
}

func New(config *model.DatabaseConfig) (*PostgresStore, errors.AppError) {
	if config == nil {
		errors.NewInternalError("postgres.storage.new_config.check.bad_arguments", "error creating storage, config is nil")
	}
	return &PostgresStore{config: config}, nil
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
func (s *PostgresStore) Database() (*sqlx.DB, errors.AppError) {
	if s.conn == nil {
		errors.NewInternalError("postgres.storage.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *PostgresStore) Open() errors.AppError {
	db, err := sqlx.Connect("pgx", s.config.Url)
	//db, err := sql.Open("pgx", s.config.Url)
	if err != nil {
		return errors.NewInternalError("postgres.storage.open.connect.fail", err.Error())
	}
	s.conn = db
	return nil
}

func (s *PostgresStore) Close() errors.AppError {
	err := s.conn.Close()
	if err != nil {
		return errors.NewInternalError("postgres.storage.close.disconnect.fail", err.Error())
	}
	s.conn = nil
	return nil

}

func ApplyFiltersToBuilder(base squirrel.SelectBuilder, filters any) squirrel.SelectBuilder {
	switch data := filters.(type) {
	case model.FilterArray:
		switch data.Connection {
		case model.AND:
			result := squirrel.And{}
			for _, bunch := range data.Filters {
				switch bunch.ConnectionType {
				case model.AND:
					lowerResult := squirrel.And{}
					for _, filter := range bunch.Bunch {
						lowerResult = append(lowerResult, applyFilter(filter))
					}
					result = append(result, lowerResult)

				case model.OR:
					lowerResult := squirrel.Or{}
					for _, filter := range bunch.Bunch {
						lowerResult = append(lowerResult, applyFilter(filter))
					}
					result = append(result, lowerResult)

				}

			}
			base = base.Where(result)
			return base
		case model.OR:
			result := squirrel.Or{}
			for _, bunch := range data.Filters {
				switch bunch.ConnectionType {
				case model.AND:
					lowerResult := squirrel.And{}
					for _, filter := range bunch.Bunch {
						lowerResult = append(lowerResult, applyFilter(filter))
					}
					result = append(result, lowerResult)
					base = base.Where(result)
					return base
				case model.OR:
					lowerResult := squirrel.Or{}
					for _, filter := range bunch.Bunch {
						lowerResult = append(lowerResult, applyFilter(filter))
					}
					result = append(result, lowerResult)
					base = base.Where(result)
					return base
				}

			}
			base = base.Where(result)
			return base
		}
	case model.FilterBunch:
		switch data.ConnectionType {
		case model.AND:
			result := squirrel.And{}
			for _, filter := range data.Bunch {
				result = append(result, applyFilter(filter))
			}

			base = base.Where(result)
			return base
		case model.OR:
			result := squirrel.Or{}
			for _, filter := range data.Bunch {
				result = append(result, applyFilter(filter))
			}
			base = base.Where(result)
			return base
		}
	}

	return base
}

func applyFilter(filter *model.Filter) squirrel.Sqlizer {
	var result squirrel.Sqlizer
	switch filter.ComparisonType {
	case model.Equal:
		result = squirrel.Eq{filter.Column: filter.Value}
	case model.GreaterThan:
		result = squirrel.Gt{filter.Column: filter.Value}
	case model.GreaterThanOrEqual:
		result = squirrel.GtOrEq{filter.Column: filter.Value}
	case model.LessThan:
		result = squirrel.Lt{filter.Column: filter.Value}
	case model.LessThanOrEqual:
		result = squirrel.LtOrEq{filter.Column: filter.Value}
	case model.NotEqual:
		result = squirrel.NotEq{filter.Column: filter.Value}
	case model.Like:
		result = squirrel.Like{filter.Column: filter.Value}
	case model.ILike:
		result = squirrel.ILike{filter.Column: filter.Value}
	}
	return result
}
