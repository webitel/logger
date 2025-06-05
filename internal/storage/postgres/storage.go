package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webitel/logger/internal/model"
	"github.com/webitel/logger/internal/storage"
	otelpgx "github.com/webitel/webitel-go-kit/infra/otel/instrumentation/pgx"
	"log/slog"

	_ "github.com/jackc/pgx/stdlib"
)

type Store struct {
	config            *model.DatabaseConfig
	conn              *pgxpool.Pool
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

func (s *Store) Database() (*pgxpool.Pool, error) {
	if s.conn == nil {
		return nil, errors.New("database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() error {
	config, err := pgxpool.ParseConfig(s.config.Url)
	if err != nil {
		return err
	}

	// Attach the OpenTelemetry tracer for pgx
	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *Store) Close() error {
	s.conn.Close()
	slog.Debug("postgres: connection closed")
	return nil

}
