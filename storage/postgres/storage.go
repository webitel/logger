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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS
	logger.object_config ( 
		id SERIAL PRIMARY KEY,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_by BIGINT NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE,
		updated_by BIGINT,
		enabled BOOLEAN NOT NULL,
		days_to_store BIGINT NOT NULL,
		period TEXT NOT NULL,
		next_upload_on TIMESTAMP WITH TIME ZONE,
		object_id BIGINT NOT NULL,
		storage_id BIGINT NOT NULL REFERENCES storage.file_backend_profiles(id),
		domain_id BIGINT NOT NULL REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE
		);`)
	if err != nil {
		return errors.NewInternalError("postgres.storage.schema_init.config_table.create", err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logger.log ( 
		id SERIAL PRIMARY KEY,
		date TIMESTAMP WITH TIME ZONE NOT NULL,
		user_id INT NOT NULL,
		user_ip TEXT NOT NULL,
		record_id BIGINT NOT NULL,
		new_state JSONB,
		action TEXT NOT NULL,
		config_id BIGINT NOT NULL REFERENCES logger.object_config(id) ON DELETE CASCADE
    );`)
	if err != nil {
		return errors.NewInternalError("postgres.storage.schema_init.log_table.create", err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logger.object_config_acl
(
    id      bigserial
        constraint object_config_acl_pk
            primary key,
    dc      bigint             not null
        constraint object_config_acl_domain_fk
            references directory.wbt_domain
            on delete cascade,
    grantor bigint
        constraint object_config_acl_grantor_id_fk
            references directory.wbt_auth
            on delete set null,
    object  integer            not null
        constraint object_config_acl_object_config_id_fk
            references logger.object_config
            on update cascade on delete cascade,
    subject bigint             not null,
    access  smallint default 0 not null,
    constraint object_config_acl_grantor_fk
        foreign key (grantor, dc) references directory.wbt_auth (id, dc)
            on update cascade on delete cascade,
    constraint object_config_acl_object_fk
        foreign key (object, dc) references logger.object_config (id, domain_id)
            on delete cascade,
    constraint object_config_acl_subject_fk
        foreign key (subject, dc) references directory.wbt_auth (id, dc)
            on delete cascade
);

alter table logger.object_config_acl
    owner to opensips;

create index IF NOT EXISTS object_config_acl_grantor_idx
    on logger.object_config_acl (grantor);

create unique index IF NOT EXISTS object_config_acl_object_subject_udx
    on logger.object_config_acl (object, subject) include (access);

create unique index IF NOT EXISTS object_config_acl_subject_object_udx
    on logger.object_config_acl (subject, object) include (access);`)
	if err != nil {
		return errors.NewInternalError("postgres.storage.schema_init.log_table.create", err.Error())
	}

	return nil
}
