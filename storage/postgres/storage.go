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

func ApplyFiltersToBuilder(base squirrel.SelectBuilder, filters ...model.Filter) squirrel.SelectBuilder {
	for _, v := range filters {
		switch v.ComparisonType {
		case 0:
			base = base.Where(squirrel.Eq{v.Column: v.Value})
		case 1:
			base = base.Where(squirrel.Gt{v.Column: v.Value})
		case 2:
			base = base.Where(squirrel.GtOrEq{v.Column: v.Value})
		case 3:
			base = base.Where(squirrel.Lt{v.Column: v.Value})
		case 4:
			base = base.Where(squirrel.LtOrEq{v.Column: v.Value})
		case 5:
			base = base.Where(squirrel.NotEq{v.Column: v.Value})
		case 6:
			base = base.Where(squirrel.Like{v.Column: v.Value})
		case 7:
			base = base.Where(squirrel.ILike{v.Column: v.Value})
		}
	}
	return base
}

//func (s *PostgresStore) SchemaInit() errors.AppError {
//	db, appErr := s.Database()
//	if appErr != nil {
//		return appErr
//	}
//	//_, err := db.Exec(`create schema if not exists logger;`)
//	//if err != nil {
//	//	return errors.NewInternalError("postgres.storage.schema_init.schema.create", err.Error())
//	//}
//	_, err := db.Exec(`
//create schema if not exists logger;
//
//
//create table if not exists logger.object_config
//(
//    id             serial
//        primary key,
//    enabled        boolean                                not null,
//    days_to_store  bigint                                 not null,
//    period         text                                   not null,
//    next_upload_on timestamp,
//    object_id      bigint                                 not null
//        constraint object_config_pk2
//            unique,
//    storage_id     bigint                                 not null
//        references storage.file_backend_profiles,
//    domain_id      bigint                                 not null
//        references directory.wbt_domain
//            on delete cascade,
//    created_at     timestamp with time zone default now() not null,
//    created_by     bigint                                 not null,
//    updated_at     timestamp with time zone,
//    updated_by     bigint,
//    constraint object_config_pk
//        unique (id, domain_id)
//);
//
//alter table logger.object_config
//    owner to opensips;
//
//create unique index if not exists object_config_id_uindex
//    on logger.object_config (id);
//
//create unique index if not exists object_config_object_id_domain_id_uindex
//    on logger.object_config (object_id, domain_id);
//
//
//
//create table if not exists log
//(
//    id        serial
//        primary key,
//    date      timestamp not null,
//    user_id   integer   not null,
//    user_ip   text      not null,
//    record_id bigint    not null,
//    new_state jsonb,
//    action    text      not null,
//    config_id integer   not null
//        references logger.object_config
//            on delete cascade
//);
//
//alter table log
//    owner to opensips;
//
//create unique index if not exists log_config_id_uindex
//    on log (config_id);
//
//create unique index if not exists log_user_id_uindex
//    on log (user_id);
//
//
//create table if not exists logger.object_config_acl
//(
//    id      bigserial
//        constraint object_config_acl_pk
//            primary key,
//    dc      bigint             not null
//        constraint object_config_acl_domain_fk
//            references directory.wbt_domain
//            on delete cascade,
//    grantor bigint
//        constraint object_config_acl_grantor_id_fk
//            references directory.wbt_auth
//            on delete set null,
//    object  integer            not null
//        constraint object_config_acl_object_config_id_fk
//            references logger.object_config
//            on update cascade on delete cascade,
//    subject bigint             not null,
//    access  smallint default 0 not null,
//    constraint object_config_acl_grantor_fk
//        foreign key (grantor, dc) references directory.wbt_auth (id, dc)
//            on update cascade on delete cascade,
//    constraint object_config_acl_object_fk
//        foreign key (object, dc) references logger.object_config (id, domain_id)
//            on delete cascade,
//    constraint object_config_acl_subject_fk
//        foreign key (subject, dc) references directory.wbt_auth (id, dc)
//            on delete cascade
//);
//
//alter table logger.object_config_acl
//    owner to opensips;
//
//create index IF NOT EXISTS object_config_acl_grantor_idx
//    on logger.object_config_acl (grantor);
//
//create unique index IF NOT EXISTS object_config_acl_object_subject_udx
//    on logger.object_config_acl (object, subject) include (access);
//
//create unique index IF NOT EXISTS object_config_acl_subject_object_udx
//    on logger.object_config_acl (subject, object) include (access);
//
//
//create function tg_obj_default_rbac() returns trigger
//    security definer
//    language plpgsql
//as
//$$
//BEGIN
//
//    EXECUTE format(
//'INSERT INTO %I.%I AS acl (dc, object, grantor, subject, access)
// SELECT $1, $2, rbac.grantor, rbac.subject, rbac.access
//   FROM (
//    -- NEW object OWNER access SUPER(255) mode (!)
//    SELECT $3, $3, (255)::int2
//     UNION ALL
//    SELECT DISTINCT ON (rbac.subject)
//      -- [WHO] grants MAX of WINDOW subset access level
//        first_value(rbac.grantor) OVER sub
//      -- [WHOM] role/user administrative unit
//      , rbac.subject
//      -- [GRANT] ALL of WINDOW subset access mode(s)
//      , bit_or(rbac.access) OVER sub
//
//      FROM directory.wbt_default_acl AS rbac
//      JOIN directory.wbt_class AS oc ON (oc.dc, oc.name) = ($1, %L)
//      -- EXISTS( OWNER membership WITH grantor role )
//      -- JOIN directory.wbt_auth_member AS sup ON (sup.role_id, sup.member_id) = (rbac.grantor, $3)
//     WHERE rbac.object = oc.id
//       AND rbac.subject <> $3
//        -- EXISTS( OWNER membership WITH grantor user/role )
//       AND (rbac.grantor = $3 OR EXISTS(SELECT true
//             FROM directory.wbt_auth_member sup
//            WHERE sup.member_id = $3
//              AND sup.role_id = rbac.grantor
//           ))
//    WINDOW sub AS (PARTITION BY rbac.subject ORDER BY rbac.access DESC)
//
//   ) AS rbac(grantor, subject, access)',
//
//--   ON CONFLICT (object, subject)
//--   DO UPDATE SET
//--     grantor = EXCLUDED.grantor,
//--     access = EXCLUDED.access',
//
//            tg_table_schema,
//            tg_table_name||'_acl',
//            tg_argv[0]::name -- objclass: directory.wbt_class.name
//        )
//    --      :srv,   :oid,   :rid
//    USING NEW.domain_id, NEW.id, NEW.created_by;
//    -- FOR EACH ROW
//    RETURN NEW;
//
//END
//$$;
//
//alter function tg_obj_default_rbac() owner to opensips;
//
//drop trigger if exists cc_agent_set_rbac_acl on "logger"."object_config";
//
//create  trigger  cc_agent_set_rbac_acl
//    after insert
//    on logger.object_config
//    for each row
//execute procedure logger.tg_obj_default_rbac('logs');
//`)
//	if err != nil {
//		return errors.NewInternalError("postgres.storage.schema_init.config_table.create", err.Error())
//	}
//
//	//	_, err = db.Exec(`create table if not exists log
//	//(
//	//    id        serial
//	//        primary key,
//	//    date      timestamp not null,
//	//    user_id   integer   not null,
//	//    user_ip   text      not null,
//	//    record_id bigint    not null,
//	//    new_state jsonb,
//	//    action    text      not null,
//	//    config_id integer   not null
//	//        references logger.object_config
//	//            on delete cascade
//	//);
//	//
//	//alter table log
//	//    owner to opensips;
//	//
//	//create unique index if not exists log_config_id_uindex
//	//    on log (config_id);
//	//
//	//create unique index if not exists log_user_id_uindex
//	//    on log (user_id);
//	//
//	//
//	//CREATE TABLE IF NOT EXISTS logger.object_config_acl
//	//(
//	//    id      bigserial
//	//        constraint object_config_acl_pk
//	//            primary key,
//	//    dc      bigint             not null
//	//        constraint object_config_acl_domain_fk
//	//            references directory.wbt_domain
//	//            on delete cascade,
//	//    grantor bigint
//	//        constraint object_config_acl_grantor_id_fk
//	//            references directory.wbt_auth
//	//            on delete set null,
//	//    object  integer            not null
//	//        constraint object_config_acl_object_config_id_fk
//	//            references logger.object_config
//	//            on update cascade on delete cascade,
//	//    subject bigint             not null,
//	//    access  smallint default 0 not null,
//	//    constraint object_config_acl_grantor_fk
//	//        foreign key (grantor, dc) references directory.wbt_auth (id, dc)
//	//            on update cascade on delete cascade,
//	//    constraint object_config_acl_object_fk
//	//        foreign key (object, dc) references logger.object_config (id, domain_id)
//	//            on delete cascade,
//	//    constraint object_config_acl_subject_fk
//	//        foreign key (subject, dc) references directory.wbt_auth (id, dc)
//	//            on delete cascade
//	//);
//	//
//	//alter table logger.object_config_acl
//	//    owner to opensips;
//	//
//	//create index IF NOT EXISTS object_config_acl_grantor_idx
//	//    on logger.object_config_acl (grantor);
//	//
//	//create unique index IF NOT EXISTS object_config_acl_object_subject_udx
//	//    on logger.object_config_acl (object, subject) include (access);
//	//
//	//create unique index IF NOT EXISTS object_config_acl_subject_object_udx
//	//    on logger.object_config_acl (subject, object) include (access);
//	//`)
//	//	if err != nil {
//	//		return errors.NewInternalError("postgres.storage.schema_init.log_table.create", err.Error())
//	//	}
//
//	//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logger.object_config_acl
//	//(
//	//    id      bigserial
//	//        constraint object_config_acl_pk
//	//            primary key,
//	//    dc      bigint             not null
//	//        constraint object_config_acl_domain_fk
//	//            references directory.wbt_domain
//	//            on delete cascade,
//	//    grantor bigint
//	//        constraint object_config_acl_grantor_id_fk
//	//            references directory.wbt_auth
//	//            on delete set null,
//	//    object  integer            not null
//	//        constraint object_config_acl_object_config_id_fk
//	//            references logger.object_config
//	//            on update cascade on delete cascade,
//	//    subject bigint             not null,
//	//    access  smallint default 0 not null,
//	//    constraint object_config_acl_grantor_fk
//	//        foreign key (grantor, dc) references directory.wbt_auth (id, dc)
//	//            on update cascade on delete cascade,
//	//    constraint object_config_acl_object_fk
//	//        foreign key (object, dc) references logger.object_config (id, domain_id)
//	//            on delete cascade,
//	//    constraint object_config_acl_subject_fk
//	//        foreign key (subject, dc) references directory.wbt_auth (id, dc)
//	//            on delete cascade
//	//);
//	//
//	//alter table logger.object_config_acl
//	//    owner to opensips;
//	//
//	//create index IF NOT EXISTS object_config_acl_grantor_idx
//	//    on logger.object_config_acl (grantor);
//	//
//	//create unique index IF NOT EXISTS object_config_acl_object_subject_udx
//	//    on logger.object_config_acl (object, subject) include (access);
//	//
//	//create unique index IF NOT EXISTS object_config_acl_subject_object_udx
//	//    on logger.object_config_acl (subject, object) include (access);
//	//
//	//`)
//	//if err != nil {
//	//	return errors.NewInternalError("postgres.storage.schema_init.log_table.create", err.Error())
//	//}
//
//	return nil
//}
