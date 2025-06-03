package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/webitel/logger/internal/storage"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/logger/internal/model"
)

type Config struct {
	storage *Store
}

var (
	configFieldsSelectMap = map[string]string{
		model.ConfigFields.Id:              "object_config.id",
		model.ConfigFields.Enabled:         "object_config.enabled",
		model.ConfigFields.DaysToStore:     "object_config.days_to_store",
		model.ConfigFields.Period:          "object_config.period",
		model.ConfigFields.NextUploadOn:    "object_config.next_upload_on",
		model.ConfigFields.DomainId:        "object_config.domain_id",
		model.ConfigFields.CreatedAt:       "object_config.created_at",
		model.ConfigFields.CreatedBy:       "author.id created_by_id, author.name created_by_name",
		model.ConfigFields.UpdatedAt:       "object_config.updated_at",
		model.ConfigFields.UpdatedBy:       "editor.id updated_by_id, editor.name updated_by_name",
		model.ConfigFields.Description:     "object_config.description",
		model.ConfigFields.LastUploadedLog: "object_config.last_uploaded_log_id",
		model.ConfigFields.LogsCount:       "(select count(l.*) from logger.log l where l.config_id = object_config.id) logs_count",
		model.ConfigFields.LogsSize:        "(select pg_size_pretty(sum(pg_column_size(l))) from logger.log l where l.config_id = object_config.id) logs_size",

		// [alias]
		model.ConfigFields.Storage: "storage.id storage_id, storage.name storage_name",
		model.ConfigFields.Object:  "object.id object_id, object.name object_name",
	}
	configFieldsFilterMap = map[string]string{
		model.ConfigFields.Id:              "object_config.id",
		model.ConfigFields.Enabled:         "object_config.enabled",
		model.ConfigFields.DaysToStore:     "object_config.days_to_store",
		model.ConfigFields.Period:          "object_config.period",
		model.ConfigFields.NextUploadOn:    "object_config.next_upload_on",
		model.ConfigFields.Object:          "object.name",
		model.ConfigFields.Storage:         "object_config.storage_id",
		model.ConfigFields.DomainId:        "object_config.domain_id",
		model.ConfigFields.CreatedAt:       "object_config.created_at",
		model.ConfigFields.CreatedBy:       "object_config.created_by",
		model.ConfigFields.UpdatedAt:       "object_config.updated_at",
		model.ConfigFields.UpdatedBy:       "object_config.updated_by",
		model.ConfigFields.Description:     "object_config.description",
		model.ConfigFields.LastUploadedLog: "object_config.last_uploaded_log_id",
	}
)

// region CONSTRUCTOR
func newConfigStore(store *Store) (storage.ConfigStore, error) {
	if store == nil {
		return nil, errors.New("error creating config interface to the config table, main store is nil")
	}
	return &Config{storage: store}, nil
}

// endregion

// region CONFIG STORAGE
func (c *Config) Update(ctx context.Context, conf *model.Config, fields []string, userId int) (*model.Config, error) {
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	if userId < 0 {
		return nil, errors.New("user id is invalid")
	}

	base := sq.Update("logger.object_config").
		Where(sq.Eq{"id": conf.Id}).
		Set("updated_at", time.Now()).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)
	if len(fields) == 0 {
		// Default all
		fields = []string{
			model.ConfigFields.Enabled,
			model.ConfigFields.DaysToStore,
			model.ConfigFields.Period,
			model.ConfigFields.NextUploadOn,
			model.ConfigFields.Storage,
			model.ConfigFields.Description,
		}
	}
	for _, v := range fields {
		switch v {
		case model.ConfigFields.Enabled:
			base = base.Set("enabled", conf.Enabled)
		case model.ConfigFields.DaysToStore:
			base = base.Set("days_to_store", conf.DaysToStore)
		case model.ConfigFields.Period:
			base = base.Set("period", conf.Period)
		case model.ConfigFields.NextUploadOn:
			base = base.Set("next_upload_on", conf.NextUploadOn)
		case model.ConfigFields.Storage:
			base = base.Set("storage_id", conf.Storage.GetId())
		case model.ConfigFields.Description:
			base = base.Set("description", conf.Description)
		case model.ConfigFields.LastUploadedLog:
			base = base.Set("last_uploaded_log_id", conf.LastUploadedLogId)
		}
	}
	query, args, _ := base.ToSql()
	query = fmt.Sprintf(
		`with p as (%s)
		select p.id,
			   p.enabled,
			   wbt_class.name             as object_name,
			   file_backend_profiles.name as storage_name,
			   p.created_at,
			   p.created_by created_by_id,
			   p.updated_at,
			   p.updated_by updated_by_id,
			   p.days_to_store,
			   p.period,
			   p.next_upload_on,
			   p.storage_id,
			   p.domain_id,
			   p.object_id,
               p.description,
               p.last_uploaded_log_id
		from p
				 LEFT JOIN directory.wbt_class ON wbt_class.id = p.object_id
				 LEFT JOIN storage.file_backend_profiles ON file_backend_profiles.id = p.storage_id`, query)
	var res model.Config
	err = pgxscan.Get(ctx, db, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Config) Insert(ctx context.Context, conf *model.Config) (*model.Config, error) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	var config model.Config
	err := pgxscan.Get(ctx, db, &config,
		`with p as (INSERT INTO
					logger.object_config (enabled, days_to_store, period, next_upload_on, storage_id, domain_id, object_id, created_by, description)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING *)
				select p.id,
					   p.enabled,
					   wbt_class.name             as object_name,
					   file_backend_profiles.name as storage_name,
					   p.created_at,
					   p.created_by created_by_id,
					   p.updated_at,
					   p.updated_by updated_by_id,
					   p.days_to_store,
					   p.period,
					   p.next_upload_on,
					   p.storage_id,
					   p.domain_id,
					   p.object_id,
					   p.description
				from p
						 LEFT JOIN directory.wbt_class ON wbt_class.id = p.object_id
						 LEFT JOIN storage.file_backend_profiles ON file_backend_profiles.id = p.storage_id`,

		conf.Enabled,
		conf.DaysToStore,
		conf.Period,
		conf.NextUploadOn,
		conf.Storage.GetId(),
		conf.DomainId,
		conf.Object.GetId(),
		conf.Author.GetId(),
		conf.Description)
	if err != nil {
		return nil, err
	}
	return &config, nil

}

func (c *Config) GetAvailableSystemObjects(ctx context.Context, domainId int, includeExisting bool, filters ...string) ([]*model.SystemObject, error) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := sq.Select("wbt_class.id id", "wbt_class.name name").
		From("directory.wbt_class").
		Where(sq.Expr("name = any(?)", pq.Array(filters))).
		Where(sq.Eq{"dc": domainId}).
		PlaceholderFormat(sq.Dollar)
	if !includeExisting {
		base = base.Where(sq.Expr(
			"id NOT IN (?)",
			sq.Select("object_id").From("logger.object_config").Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId}),
		))
	}
	query, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}
	var objects []*model.SystemObject
	err = pgxscan.Select(ctx, db, &objects, query, args...)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func (c *Config) CheckAccess(ctx context.Context, domainId, id int, groups []int, access uint8) (bool, error) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return false, appErr
	}
	subquery := sq.Select("1").From("logger.object_config_acl acl").
		Where(sq.Eq{"acl.dc": domainId}).
		Where(sq.Eq{"acl.object": id}).
		Where("acl.subject = any( ?::int[])", pq.Array(groups)).
		Where("acl.access & ? = ?", access, access).PlaceholderFormat(sq.Dollar)
	base := sq.Select("1").Where(sq.Expr("exists(?)", subquery))
	query, args, err := base.ToSql()
	if err != nil {
		return false, err
	}
	var permission bool
	err = pgxscan.Get(ctx, db, &permission, query, args...)
	if err != nil {
		return false, err
	}
	return permission, nil
}

func (c *Config) Delete(ctx context.Context, id int, domainId int) (int, error) {
	db, err := c.storage.Database()
	if err != nil {
		return 0, err
	}
	base := sq.Delete("logger.object_config").Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.Id]: id}).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId}).PlaceholderFormat(sq.Dollar)
	query, args, err := base.ToSql()
	if err != nil {
		return 0, err
	}
	tag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), err
}

func (c *Config) DeleteMany(ctx context.Context, rbac *model.RbacOptions, ids []int, domainId int) (int, error) {
	db, err := c.storage.Database()
	if err != nil {
		return 0, err
	}
	base := sq.Delete("logger.object_config").Where(sq.Expr(configFieldsFilterMap[model.ConfigFields.Id]+" = any(?::int[])", pq.Array(ids))).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId}).PlaceholderFormat(sq.Dollar)
	if rbac != nil {
		subquery := sq.Select("1").From("logger.object_config_acl acl").
			Where("acl.dc = "+configFieldsFilterMap[model.ConfigFields.DomainId]).
			Where("acl.object = "+configFieldsFilterMap[model.ConfigFields.DomainId]).
			Where("acl.subject = any( ?::int[])", pq.Array(rbac.Groups)).
			Where("acl.access & ? = ?", rbac.Access, rbac.Access).
			Limit(1)
		base = base.Where(sq.Expr("exists(?)", subquery))
	}
	query, args, err := base.ToSql()
	if err != nil {
		return 0, err
	}
	tag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), err
}

func (c *Config) GetByObjectId(ctx context.Context, domainId int, objectId int) (*model.Config, error) {
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	base := c.GetQueryBase(c.getFields(), nil).Where(
		sq.Eq{configFieldsFilterMap[model.ConfigFields.Object]: objectId},
		sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId},
	)
	var config model.Config
	query, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(ctx, db, &config, query, args...)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Get(ctx context.Context, rbac *model.RbacOptions, id int, domainId int) (*model.Config, error) {
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	base := c.GetQueryBase(c.getFields(), rbac).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.Id]: id}).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId})
	var config model.Config
	query, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(ctx, db, &config, query, args...)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Select(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, filters any) ([]*model.Config, error) {
	var (
		sql  string
		args []any
	)
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	base, err := storage.ApplyFiltersToBuilderBulk(c.GetQueryBaseFromSearchOptions(opt, rbac), configFieldsFilterMap, filters)
	if err != nil {
		return nil, err
	}
	switch req := base.(type) {
	case sq.SelectBuilder:
		sql, args, _ = req.ToSql()
	default:
		return nil, errors.New("wrong base type")
	}
	var configs []*model.Config
	err = pgxscan.Select(ctx, db, &configs, sql, args...)
	if err != nil {
		return nil, err
	}
	return configs, nil

}

// endregion

// region SYSTEM FUNCTIONS

func (c *Config) GetQueryBaseFromSearchOptions(opt *model.SearchOptions, rbac *model.RbacOptions) sq.SelectBuilder {
	var fields []string

	if opt == nil {
		return c.GetQueryBase(c.getFields(), rbac)
	}
	for _, v := range opt.Fields {
		if columnName, ok := configFieldsSelectMap[v]; ok {
			fields = append(fields, columnName)
		} else {
			fields = append(fields, v)
		}
	}
	if len(fields) == 0 {
		fields = append(fields,
			c.getFields()...)
	}
	base := c.GetQueryBase(fields, rbac)
	if opt.Search != "" {
		base = base.Where(sq.ILike{"object.name": opt.Search + "%"})
	}
	if opt.Sort != "" {
		splitted := strings.Split(opt.Sort, ":")
		if len(splitted) == 2 {
			order := splitted[0]
			column := splitted[1]
			// Lookup columns -- order by name
			switch column {
			case model.ConfigFields.Object:
				column = "object_name"
			case model.ConfigFields.Storage:
				column = "storage_name"
			}
			base = base.OrderBy(fmt.Sprintf("%s %s", column, order))
		}
	}
	if opt.Size > 0 {
		offset := (opt.Page - 1) * opt.Size
		if offset < 0 {
			offset = 0
		}
		base = base.Offset(uint64(offset)).Limit(uint64(opt.Size))
	}

	return base
}

func (c *Config) GetQueryBase(fields []string, rbac *model.RbacOptions) sq.SelectBuilder {
	base := sq.Select(fields...).From("logger.object_config").
		LeftJoin("directory.wbt_class object ON object.id = object_config.object_id").
		LeftJoin("storage.file_backend_profiles storage ON storage.id = object_config.storage_id").
		LeftJoin("directory.wbt_user author ON author.id = object_config.created_by").
		LeftJoin("directory.wbt_user editor ON editor.id = object_config.updated_by").
		PlaceholderFormat(sq.Dollar)
	return c.insertRbacCondition(base, rbac)
}

func (c *Config) insertRbacCondition(base sq.SelectBuilder, rbac *model.RbacOptions) sq.SelectBuilder {
	if rbac != nil {
		subquery := sq.Select("1").From("logger.object_config_acl acl").
			Where("acl.dc = object_config.domain_id").
			Where("acl.object = object_config.id").
			Where("acl.subject = any( ?::int[])", pq.Array(rbac.Groups)).
			Where("acl.access & ? = ?", rbac.Access, rbac.Access).
			Limit(1)
		base = base.Where(sq.Expr("exists(?)", subquery))
	}
	return base
}

func (c *Config) getFields() []string {
	var fields []string
	for _, value := range configFieldsSelectMap {
		fields = append(fields, value)
	}
	return fields
}

// endregion
