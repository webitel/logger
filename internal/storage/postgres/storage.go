package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/webitel/logger/internal/model"
	"github.com/webitel/logger/internal/storage"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"log/slog"

	_ "github.com/jackc/pgx/stdlib"
)

type Store struct {
	config            *model.DatabaseConfig
	conn              *sqlx.DB
	logStore          storage.LogStore
	configStore       storage.ConfigStore
	loginAttemptStore storage.LoginAttemptStore
}

func New(config *model.DatabaseConfig) *Store {

	return &Store{config: config}
}

func (s *Store) Log() storage.LogStore {
	if s.logStore == nil {
		log, err := newLogStore(s)
		if err != nil {
			return nil
		}
		s.logStore = log
	}
	return s.logStore
}
func (s *Store) Config() storage.ConfigStore {
	if s.configStore == nil {
		conf, err := newConfigStore(s)
		if err != nil {
			return nil
		}
		s.configStore = conf
	}
	return s.configStore
}
func (s *Store) LoginAttempt() storage.LoginAttemptStore {
	if s.loginAttemptStore == nil {
		conf, err := newLoginAttemptStore(s)
		if err != nil {
			return nil
		}
		s.loginAttemptStore = conf
	}
	return s.loginAttemptStore
}

func (s *Store) Database() (*sqlx.DB, error) {
	if s.conn == nil {
		return nil, model.NewInternalError("postgres.storage.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() error {
	driver := "pgx"
	db, err := otelsql.Open(driver, s.config.Url, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL, attribute.String("driver", driver),
	))
	if err != nil {
		return model.NewInternalError("postgres.storage.open.connect.fail", err.Error())
	}
	s.conn = sqlx.NewDb(db, driver)
	slog.Debug(fmt.Sprintf("postgres: connection opened"))
	return nil
}

func (s *Store) Close() error {
	err := s.conn.Close()
	if err != nil {
		return model.NewInternalError("postgres.storage.close.disconnect.fail", fmt.Sprintf("postgres: %s", err.Error()))
	}
	s.conn = nil
	slog.Debug(fmt.Sprintf("postgres: connection closed"))
	return nil

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
