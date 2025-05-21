package postgres

import (
	"errors"
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
		return nil, errors.New("database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() error {
	driver := "pgx"
	db, err := otelsql.Open(driver, s.config.Url, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL, attribute.String("driver", driver),
	))
	if err != nil {
		return err
	}
	s.conn = sqlx.NewDb(db, driver)
	slog.Debug("postgres: connection opened")
	return nil
}

func (s *Store) Close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	s.conn = nil
	slog.Debug("postgres: connection closed")
	return nil

}
