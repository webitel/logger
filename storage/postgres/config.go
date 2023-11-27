package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/storage"
	"github.com/webitel/wlog"

	errors "github.com/webitel/engine/model"
)

type Config struct {
	storage storage.Storage
}

var (
	configFieldsMap = map[string]string{
		"id":                   "object_config.id",
		"enabled":              "object_config.enabled",
		"days_to_store":        "object_config.days_to_store",
		"period":               "object_config.period",
		"next_upload_on":       "object_config.next_upload_on",
		"object_name":          "wbt_class.name as object_name",
		"object_id":            "object_config.object_id",
		"storage_id":           "object_config.storage_id",
		"storage_name":         "file_backend_profiles.name as storage_name",
		"domain_id":            "object_config.domain_id",
		"created_at":           "object_config.created_at",
		"created_by":           "object_config.created_by",
		"updated_at":           "object_config.updated_at",
		"updated_by":           "object_config.updated_by",
		"description":          "object_config.description",
		"last_uploaded_log_id": "object_config.last_uploaded_log_id",

		// [alias]
		"storage": "object_config.storage_id, file_backend_profiles.name as storage_name",
		"object":  "object_config.object_id, wbt_class.name as object_name",
	}
)

// region CONSTRUCTOR
func newConfigStore(store storage.Storage) (storage.ConfigStore, errors.AppError) {
	if store == nil {
		return nil, errors.NewInternalError("postgres.config.new_config.check.bad_arguments", "error creating config interface to the config table, main store is nil")
	}
	return &Config{storage: store}, nil
}

// endregion

// region CONFIG STORAGE
func (c *Config) Update(ctx context.Context, conf *model.Config, fields []string, userId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	if userId < 0 {
		return nil, errors.NewBadRequestError("postgres.config.update.invalid_args.user_id", "user id is invalid !")
	}

	base := sq.Update("logger.object_config").
		Where(sq.Eq{"id": conf.Id}).
		Set("updated_at", time.Now()).
		PlaceholderFormat(sq.Dollar)
	if len(fields) == 0 {
		fields = []string{"enabled", "days_to_store", "period", "next_upload_on", "storage", "description"}
	}
	for _, v := range fields {
		switch v {
		case "enabled":
			base = base.Set("enabled", conf.Enabled)
		case "days_to_store":
			base = base.Set("days_to_store", conf.DaysToStore)
		case "period":
			base = base.Set("period", conf.Period)
		case "next_upload_on":
			base = base.Set("next_upload_on", conf.NextUploadOn)
		case "storage":
			base = base.Set("storage_id", conf.Storage.Id)
		case "description":
			base = base.Set("description", conf.Description)
		case "last_uploaded_log_id":
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
		return nil, errors.NewInternalError("postgres.config.update.query.fail", err.Error())
	}
	defer res.Close()
	row, appErr := c.ScanRow(res)
	if appErr != nil {
		return nil, appErr
	}
	return row, nil
}

func (c *Config) Insert(ctx context.Context, conf *model.Config, userId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	res, err := db.QueryContext(ctx,
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
		return nil, errors.NewInternalError("postgres.config.insert.query.error", err.Error())
	}
	defer res.Close()
	row, appErr := c.ScanRow(res)
	if appErr != nil {
		return nil, appErr
	}
	return row, nil
}

func (c *Config) GetAvailableSystemObjects(ctx context.Context, domainId int, includeExisting bool, filters ...string) ([]*model.Lookup, errors.AppError) {
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
			sq.Select("object_id").From("logger.object_config").Where(sq.Eq{"domain_id": domainId}),
		))
	}
	// endregion
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	// region PREFORM
	rows, err := base.RunWith(db).QueryContext(ctx)

	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_available_objects.query.fail", err.Error())
	}
	defer rows.Close()
	// endregion
	// region SCAN
	var res []*model.Lookup
	for rows.Next() {
		var obj model.Lookup
		r := rows.Scan(&obj.Id, &obj.Name)
		if r != nil {
			return nil, errors.NewInternalError("postgres.config.get_available_objects.scan.fail", r.Error())
		}
		res = append(res, &obj)
	}
	if appErr != nil {
		return nil, appErr
	}
	// endregion
	return res, nil
}

func (c *Config) CheckAccess(ctx context.Context, domainId, id int64, groups []int, access uint32) (bool, errors.AppError) {
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
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	res := base.RunWith(db).QueryRowContext(ctx)
	var ac bool
	err := res.Scan(&ac)
	if err != nil {
		return false, errors.NewInternalError("postgres.config.check_access.scan.error", err.Error())
	}
	return ac, nil
}

func (c *Config) Delete(ctx context.Context, id int32) errors.AppError {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return appErr
	}
	base := sq.Delete("logger.object_config").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return errors.NewInternalError("postgres.config.delete.query.error", err.Error())
	}
	if i, err := res.RowsAffected(); err == nil && i == 0 {
		return errors.NewBadRequestError("postgres.config.delete.result.no_rows_for_delete", err.Error())
	}
	return nil
}

func (c *Config) DeleteMany(ctx context.Context, rbac *model.RbacOptions, ids []int32) errors.AppError {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return appErr
	}
	base := sq.Delete("logger.object_config").Where(sq.Expr("object_config.id = any(?::int[])", pq.Array(ids))).PlaceholderFormat(sq.Dollar)
	if rbac != nil {
		subquery := sq.Select("1").From("logger.object_config_acl acl").
			Where("acl.dc = object_config.domain_id").
			Where("acl.object = object_config.id").
			Where("acl.subject = any( ?::int[])", pq.Array(rbac.Groups)).
			Where("acl.access & ? = ?", rbac.Access, rbac.Access).
			Limit(1)
		base = base.Where(sq.Expr("exists(?)", subquery))
	}
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return errors.NewInternalError("postgres.config.delete_many.query.error", err.Error())
	}
	if i, err := res.RowsAffected(); err == nil && i == 0 {
		return errors.NewBadRequestError("postgres.config.delete_many.result.no_rows_for_delete", "no rows were affected while deleting")
	}
	return nil
}

func (c *Config) GetByObjectId(ctx context.Context, domainId int, objId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBase(c.getFields(), nil).Where(sq.Eq{"object_config.object_id": objId}, sq.Eq{"object_config.domain_id": domainId})
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_object.query.fail", err.Error())
	}
	defer rows.Close()
	configs, appErr := c.ScanRow(rows)
	if appErr != nil {
		return nil, appErr
	}
	return configs, nil
}

func (c *Config) GetById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBase(c.getFields(), rbac).Where(sq.Eq{"object_config.id": id})
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_id.query.fail", err.Error())
	}
	defer rows.Close()
	config, appErr := c.ScanRow(rows)
	if appErr != nil {
		return nil, appErr
	}
	return config, nil
}

func (c *Config) Get(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, filters any) ([]*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := ApplyFiltersToBuilder(c.GetQueryBaseFromSearchOptions(opt, rbac), configFieldsMap, filters)
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	rows, err := base.RunWith(db).QueryContext(ctx)

	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get.query_execute.fail", err.Error())
	}
	defer rows.Close()
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

// endregion

// region SYSTEM FUNCTIONS

func (c *Config) ScanRow(rows *sql.Rows) (*model.Config, errors.AppError) {
	res, err := c.ScanRows(rows)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func (c *Config) ScanRows(rows *sql.Rows) ([]*model.Config, errors.AppError) {
	if rows == nil {
		return nil, errors.NewInternalError("postgres.config.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.scan.get_columns.error", err.Error())
	}
	var configs []*model.Config

	for rows.Next() {
		var config model.Config
		binds := make([]func(dst *model.Config) interface{}, 0, 0)
		for _, v := range cols {

			switch v {
			case "id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Id })
			case "enabled":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Enabled })
			case "days_to_store":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.DaysToStore })
			case "period":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Period })
			case "next_upload_on":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.NextUploadOn })
			case "object_id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Object.Id })
			case "object_name":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Object.Name })
			case "storage_id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Storage.Id })
			case "storage_name":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Storage.Name })
			case "domain_id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.DomainId })
			case "created_at":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.CreatedAt })
			case "created_by":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.CreatedBy })
			case "updated_at":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.UpdatedAt })
			case "updated_by":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.UpdatedBy })
			case "description":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.Description })
			case "last_uploaded_log_id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.LastUploadedLog })
			default:
				panic("postgres.log.scan.get_columns.error: columns gotten from sql don't respond to model columns. Unknown column: " + v)
			}

		}
		bindFunc := func(binds []func(config *model.Config) any) []any {
			var fields []any
			for _, v := range binds {
				fields = append(fields, v(&config))
			}
			return fields
		}
		err = rows.Scan(bindFunc(binds)...)
		if err != nil {
			return nil, errors.NewInternalError("postgres.config.scan.scan.error", err.Error())
		}
		configs = append(configs, &config)
	}
	if len(configs) == 0 {
		return nil, errors.NewBadRequestError("postgres.config.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return configs, nil
}

func (c *Config) GetQueryBaseFromSearchOptions(opt *model.SearchOptions, rbac *model.RbacOptions) sq.SelectBuilder {
	var fields []string

	if opt == nil {
		return c.GetQueryBase(c.getFields(), rbac)
	}
	for _, v := range opt.Fields {
		fields = append(fields, configFieldsMap[v])
		if columnName, ok := logFieldsMap[v]; ok {
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
			if column == "object" {
				column = "object_name"
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
	base := sq.Select(fields...).From("logger.object_config").JoinClause("LEFT JOIN directory.wbt_class ON wbt_class.id = object_config.object_id").JoinClause("LEFT JOIN storage.file_backend_profiles ON file_backend_profiles.id = object_config.storage_id").PlaceholderFormat(sq.Dollar)
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
	for _, value := range configFieldsMap {
		fields = append(fields, value)
	}
	return fields
}

// endregion
