package postgres

import (
	"database/sql"
	"webitel_logger/model"
	"webitel_logger/storage"

	_ "github.com/jackc/pgx/stdlib"
	errors "github.com/webitel/engine/model"
)

type PostgresStore struct {
	config      *model.DatabaseConfig
	conn        *sql.DB
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
	return nil
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
func (s *PostgresStore) Database() (*sql.DB, errors.AppError) {
	if s.conn == nil {
		errors.NewInternalError("postgres.storage.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *PostgresStore) Open() errors.AppError {
	db, err := sql.Open("pgx", s.config.Url)
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

func (s *PostgresStore) SchemaInit() errors.AppError {
	db, appErr := s.Database()
	if appErr != nil {
		return appErr
	}
	_, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS logger;`)
	if err != nil {
		return errors.NewInternalError("postgres.storage.schema_init.schema.create", err.Error())
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logger.log ( 
		id SERIAL PRIMARY KEY,
		date TIMESTAMP NOT NULL,
		user_id TEXT NOT NULL,
		user_ip TEXT NOT NULL,
		object_id BIGINT REFERENCES directory.wbt_class(id),
		new_state TEXT,
		domain_id BIGINT NOT NULL);`)
	if err != nil {
		return errors.NewInternalError("postgres.storage.schema_init.log_table.create", err.Error())
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS
	logger.object_config ( 
		id SERIAL PRIMARY KEY,
		enabled BOOLEAN NOT NULL,
		days_to_store BIGINT NOT NULL,
		period TEXT NOT NULL,
		next_upload_on TIMESTAMP,
		object_id BIGINT REFERENCES directory.wbt_class(id),
		storage_id BIGINT REFERENCES storage.file_backend_profiles(id),
		domain_id INT NOT NULL);`)
	if err != nil {
		return errors.NewInternalError("postgres.storage.schema_init.config_table.create", err.Error())
	}
	return nil
}
