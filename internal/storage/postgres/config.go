package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/webitel/logger/internal/storage"
	"github.com/webitel/logger/internal/storage/postgres/utils"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/logger/internal/model"
)

type Config struct {
	storage storage.Storage
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
		model.ConfigFields.CreatedBy:       "object_config.created_by",
		model.ConfigFields.UpdatedAt:       "object_config.updated_at",
		model.ConfigFields.UpdatedBy:       "object_config.updated_by",
		model.ConfigFields.Description:     "object_config.description",
		model.ConfigFields.LastUploadedLog: "object_config.last_uploaded_log_id",
		model.ConfigFields.LogsCount:       "(select count(l.*) from logger.log l where l.config_id = object_config.id) logs_count",
		model.ConfigFields.LogsSize:        "(select pg_size_pretty(sum(pg_column_size(l))) from logger.log l where l.config_id = object_config.id) logs_size",

		// [alias]
		model.ConfigFields.Storage: "object_config.storage_id, file_backend_profiles.name as storage_name",
		model.ConfigFields.Object:  "object_config.object_id, wbt_class.name as object_name",
	}
	configFieldsFilterMap = map[string]string{
		model.ConfigFields.Id:              "object_config.id",
		model.ConfigFields.Enabled:         "object_config.enabled",
		model.ConfigFields.DaysToStore:     "object_config.days_to_store",
		model.ConfigFields.Period:          "object_config.period",
		model.ConfigFields.NextUploadOn:    "object_config.next_upload_on",
		model.ConfigFields.Object:          "object_config.object_id",
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
func newConfigStore(store storage.Storage) (storage.ConfigStore, error) {
	if store == nil {
		return nil, errors.New("error creating config interface to the config table, main store is nil")
	}
	return &Config{storage: store}, nil
}

// endregion

// region CONFIG STORAGE
func (c *Config) Update(ctx context.Context, conf *model.Config, fields []string, userId int64) (*model.Config, error) {
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
			base = base.Set("storage_id", conf.Storage.Id)
		case model.ConfigFields.Description:
			base = base.Set("description", conf.Description)
		case model.ConfigFields.LastUploadedLog:
			base = base.Set("last_uploaded_log_id", conf.LastUploadedLog)
		}
	}
	query, args, _ := base.ToSql()
	query = fmt.Sprintf(
		`with p as (%s RETURNING *)
		select p.id,
			   p.enabled,
			   wbt_class.name             as object_name,
			   file_backend_profiles.name as storage_name,
			   p.created_at,
			   p.created_by,
			   p.updated_at,
			   p.updated_by,
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
	res, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return utils.ScanRow(res, c.GetScanPlan)
}

func (c *Config) Insert(ctx context.Context, conf *model.Config, userId int64) (*model.Config, error) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	row, err := db.QueryContext(ctx,
		`with p as (INSERT INTO
					logger.object_config (enabled, days_to_store, period, next_upload_on, storage_id, domain_id, object_id, created_by, description)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING *)
				select p.id,
					   p.enabled,
					   wbt_class.name             as object_name,
					   file_backend_profiles.name as storage_name,
					   p.created_at,
					   p.created_by,
					   p.updated_at,
					   p.updated_by,
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
		conf.Storage.Id,
		conf.DomainId,
		conf.Object.Id,
		userId,
		conf.Description)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	return utils.ScanRow(row, c.GetScanPlan)
}

func (c *Config) GetAvailableSystemObjects(ctx context.Context, domainId int, includeExisting bool, filters ...string) ([]*model.Lookup, error) {
	// region CREATING QUERY
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := sq.Select("wbt_class.id", "wbt_class.name").
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
	// endregion
	// region PREFORM
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// endregion
	// region SCAN
	var res []*model.Lookup
	for rows.Next() {
		var obj model.Lookup
		err := rows.Scan(&obj.Id, &obj.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, &obj)
	}
	// endregion
	return res, nil
}

func (c *Config) CheckAccess(ctx context.Context, domainId, id int64, groups []int64, access uint8) (bool, error) {
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
	res := base.RunWith(db).QueryRowContext(ctx)
	var ac bool
	err := res.Scan(&ac)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return ac, nil
}

func (c *Config) Delete(ctx context.Context, id int32, domainId int64) (int, error) {
	db, err := c.storage.Database()
	if err != nil {
		return 0, err
	}
	base := sq.Delete("logger.object_config").Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.Id]: id}).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId}).PlaceholderFormat(sq.Dollar)
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	return int(affected), err
}

func (c *Config) DeleteMany(ctx context.Context, rbac *model.RbacOptions, ids []int32, domainId int64) (int, error) {
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
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return 0, err
	}
	i, err := res.RowsAffected()
	return int(i), err
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
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return utils.ScanRow(rows, c.GetScanPlan)
}

func (c *Config) GetById(ctx context.Context, rbac *model.RbacOptions, id int, domainId int64) (*model.Config, error) {
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	base := c.GetQueryBase(c.getFields(), rbac).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.Id]: id}).Where(sq.Eq{configFieldsFilterMap[model.ConfigFields.DomainId]: domainId})
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return utils.ScanRow(rows, c.GetScanPlan)
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
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return utils.ScanRows(rows, c.GetScanPlan)

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
		base = base.Where(sq.ILike{"wbt_class.name": opt.Search + "%"})
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
	offset := (opt.Page - 1) * opt.Size
	if offset < 0 {
		offset = 0
	}
	return base.Offset(uint64(offset)).Limit(uint64(opt.Size))
}

func (c *Config) GetQueryBase(fields []string, rbac *model.RbacOptions) sq.SelectBuilder {
	base := sq.Select(fields...).From("logger.object_config").
		JoinClause("LEFT JOIN directory.wbt_class ON wbt_class.id = object_config.object_id").
		JoinClause("LEFT JOIN storage.file_backend_profiles ON file_backend_profiles.id = object_config.storage_id").
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

func (c *Config) GetScanPlan(columns []string) []func(*model.Config) any {
	var binds []func(*model.Config) any
	for _, v := range columns {
		var bind func(*model.Config) any
		switch v {
		case "id":
			bind = func(dst *model.Config) interface{} { return &dst.Id }
		case "enabled":
			bind = func(dst *model.Config) interface{} { return &dst.Enabled }
		case "days_to_store":
			bind = func(dst *model.Config) interface{} { return &dst.DaysToStore }
		case "period":
			bind = func(dst *model.Config) interface{} { return &dst.Period }
		case "next_upload_on":
			bind = func(dst *model.Config) interface{} { return &dst.NextUploadOn }
		case "object_id":
			bind = func(dst *model.Config) interface{} { return &dst.Object.Id }
		case "object_name":
			bind = func(dst *model.Config) interface{} { return &dst.Object.Name }
		case "storage_id":
			bind = func(dst *model.Config) interface{} { return &dst.Storage.Id }
		case "storage_name":
			bind = func(dst *model.Config) interface{} { return &dst.Storage.Name }
		case "domain_id":
			bind = func(dst *model.Config) interface{} { return &dst.DomainId }
		case "created_at":
			bind = func(dst *model.Config) interface{} { return &dst.CreatedAt }
		case "created_by":
			bind = func(dst *model.Config) interface{} { return &dst.CreatedBy }
		case "updated_at":
			bind = func(dst *model.Config) interface{} { return &dst.UpdatedAt }
		case "updated_by":
			bind = func(dst *model.Config) interface{} { return &dst.UpdatedBy }
		case "description":
			bind = func(dst *model.Config) interface{} { return &dst.Description }
		case "last_uploaded_log_id":
			bind = func(dst *model.Config) interface{} { return &dst.LastUploadedLog }
		case "logs_count":
			bind = func(dst *model.Config) interface{} { return &dst.LogsCount }
		case "logs_size":
			bind = func(dst *model.Config) interface{} { return &dst.LogsSize }
		default:
			bind = func(dst *model.Config) any { return nil }
		}
		binds = append(binds, bind)
	}
	return binds
}

// endregion
