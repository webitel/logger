package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strings"
	"time"
	"webitel_logger/model"
	"webitel_logger/storage"

	errors "github.com/webitel/engine/model"
)

type Config struct {
	storage storage.Storage
}

func newConfigStore(store storage.Storage) (storage.ConfigStore, errors.AppError) {
	if store == nil {
		return nil, errors.NewInternalError("postgres.config.new_config.check.bad_arguments", "error creating config interface to the config table, main store is nil")
	}
	return &Config{storage: store}, nil
}

func (c *Config) Update(ctx context.Context, conf *model.Config, fields []string, userId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := sq.Update("logger.object_config").Where(sq.Eq{"object_config.id": conf.Id}).Set("updated_at", time.Now()).Set("updated_by", userId).PlaceholderFormat(sq.Dollar)
	if len(fields) == 0 {
		fields = []string{"enabled", "days_to_store", "period", "next_upload_on", "storage"}
	}
	for _, v := range fields {
		switch v {
		case "enabled":
			base = base.Set("object_config.enabled", conf.Enabled)
		case "days_to_store":
			base = base.Set("object_config.days_to_store", conf.DaysToStore)
		case "period":
			base = base.Set("object_config.period", conf.Period)
		case "next_upload_on":
			base = base.Set("object_config.next_upload_on", conf.NextUploadOn)
		case "storage":
			base = base.Set("object_config.storage_id", conf.Storage.Id)
		}
	}

	//row, err := db.QueryContext(ctx,
	//	`UPDATE
	//		logger.object_config
	//	SET
	//		enabled = $1,
	//		days_to_store = $2,
	//		period = $3,
	//		next_upload_on = $4,
	//		storage_id = $5,
	//		updated_at = $6,
	//		updated_by = $7
	//	WHERE
	//		id = $8
	//	RETURNING id, enabled, created_at, created_by, updated_at, updated_by, days_to_store, period, next_upload_on, storage_id, domain_id, object_id;`,
	//	conf.Enabled, conf.DaysToStore, conf.Period, conf.NextUploadOn, conf.Storage.Id, time.Now(), userId, conf.Id)
	row, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.update.query.fail", err.Error())
	}
	newConfig, appErr := c.ScanRow(row)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.update.scan.fail", err.Error())
	}
	return newConfig, nil
}

//func (c *Config) CheckAccess(ctx context.Context, domainId, id int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError) {
//	var i int
//	db, appErr := c.storage.Database()
//	if appErr != nil {
//		return false, appErr
//	}
//	res := db.QueryRowContext(ctx, `select 1
//		where exists(
//          select 1
//          from call_center.cc_queue_acl a
//          where a.dc = $1
//            and a.object = $2
//            and a.subject = any ($3::int[])
//            and a.access & $4 = $4
//        )`, domainId, id, groups, access)
//	err := res.Scan(&i)
//	if err != nil {
//		return false, nil
//	}
//
//	return i == 1, nil
//}

//func (c *Config) CheckAccessByObjectId(ctx context.Context, domainId, objectId int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError) {
//	var i int
//	db, appErr := c.storage.Database()
//	if appErr != nil {
//		return false, appErr
//	}
//	res := db.QueryRowContext(ctx, `select 1
//		where exists(
//          select 1
//          from call_center.cc_queue_acl a
//          where a.dc = $1
//            and a.object IN (select id from logger.object_config where object_config.domain_id = $1 and object_config.object_id = $2 )
//            and a.subject = any ($3::int[])
//            and a.access & $4 = $4
//        )`, domainId, objectId, groups, access)
//	err := res.Scan(&i)
//	if err != nil {
//		return false, nil
//	}
//
//	return i == 1, nil
//}

func (c *Config) Insert(ctx context.Context, conf *model.Config, userId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	res, err := db.QueryContext(ctx,
		`INSERT INTO
		   logger.object_config(enabled, days_to_store, period, next_upload_on, storage_id, domain_id, object_id, created_by) 
		VALUES
		   (
			  $1, $2, $3, $4, $5, $6, $7, $8 
		   )
		   RETURNING id, enabled, created_at, created_by, updated_at, updated_by, days_to_store, period, next_upload_on, storage_id, domain_id, object_id`,
		conf.Enabled,
		conf.DaysToStore,
		conf.Period,
		conf.NextUploadOn,
		conf.Storage.Id,
		conf.DomainId,
		conf.Object.Id,
		userId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.insert.query.fail", err.Error())
	}
	newConfig, appErr := c.ScanRow(res)
	if appErr != nil {
		return nil, appErr
	}
	return newConfig, nil
}

func (c *Config) GetByObjectId(ctx context.Context, domainId int, objId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBase(c.getFields(), nil).Where(sq.Eq{"object_config.object_id": objId}, sq.Eq{"object_config.domain_id": domainId})

	rows, err := base.RunWith(db).QueryContext(ctx)

	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_object.query.fail", err.Error())
	}
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
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_id.query.fail", err.Error())
	}
	config, appErr := c.ScanRow(rows)
	if appErr != nil {
		return nil, appErr
	}
	return config, nil
}

//func (c *Config) GetAllRbac(ctx context.Context, opt *model.SearchOptions, domainId int, groups []int, access auth_manager.PermissionAccess) (*[]model.Config, errors.AppError) {
//	db, appErr := c.storage.Database()
//	if appErr != nil {
//		return nil, appErr
//	}
//	base := c.GetQueryBaseFromSearchOptions(opt)
//
//	rows, err := base.Where(sq.Eq{"domain_id": domainId}, sq.Expr("exists(select 1 from logger.object_config_acl acl where acl.dc = object_config.domain_id and acl.object = object_config.id and acl.subject = any($1::int[]) and acl.access&$2 = $2)"), groups, access).RunWith(db).QueryContext(ctx)
//	if err != nil {
//		return nil, errors.NewInternalError("postgres.config.get_all.query.fail", err.Error())
//	}
//	configs, appErr := c.ScanRows(rows)
//	if appErr != nil {
//		return nil, appErr
//	}
//	return &configs, nil
//}

func (c *Config) Get(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, filters ...model.Filter) (*[]model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := ApplyFiltersToBuilder(c.GetQueryBaseFromSearchOptions(opt, rbac), filters...)
	sql, _, _ := base.ToSql()
	fmt.Println(sql)
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get.query_execute.fail", err.Error())
	}
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return &res, nil
}

func (c *Config) ScanRow(rows *sql.Rows) (*model.Config, errors.AppError) {
	res, err := c.ScanRows(rows)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (c *Config) ScanRows(rows *sql.Rows) ([]model.Config, errors.AppError) {
	if rows == nil {
		return nil, errors.NewInternalError("postgres.config.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.scan.get_columns.error", err.Error())
	}
	var configs []model.Config

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
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, errors.NewBadRequestError("postgres.config.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return configs, nil
}

// Refactor GetQueryBaseFromSearchOptions for beauty

func (c *Config) GetQueryBaseFromSearchOptions(opt *model.SearchOptions, rbac *model.RbacOptions) sq.SelectBuilder {
	var fields []string
	if opt == nil {
		return c.GetQueryBase(c.getFields(), rbac)
	}
	for _, v := range opt.Fields {
		switch v {
		case "id":
			fields = append(fields, "object_config.id")
		case "enabled":
			fields = append(fields, "object_config.enabled")
		case "days_to_store":
			fields = append(fields, "object_config.days_to_store")
		case "period":
			fields = append(fields, "object_config.period")
		case "next_upload_on":
			fields = append(fields, "object_config.next_upload_on")
		case "object":
			fields = append(fields, "object_config.object_id")
			fields = append(fields, "wbt_class.name as object_name")
		case "storage":
			fields = append(fields, "object_config.storage_id")
			fields = append(fields, "file_backend_profiles.name as storage_name")
		case "domain_id":
			fields = append(fields, "object_config.domain_id")
		case "created_at":
			fields = append(fields, "object_config.created_at")
		case "created_by":
			fields = append(fields, "object_config.created_by")
		case "updated_at":
			fields = append(fields, "object_config.updated_at")
		case "updated_by":
			fields = append(fields, "object_config.updated_by")
		}
	}
	if len(fields) == 0 {
		fields = append(fields,
			c.getFields()...)
	}
	base := c.GetQueryBase(fields, rbac)
	if opt.Search != "" {
		base = base.Where(sq.Like{"wbt_class.name": "%" + strings.ToLower(opt.Search) + "%"})
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
		base = base.Where(sq.Expr("exists(select 1 from logger.object_config_acl acl where acl.dc = object_config.domain_id and acl.object = object_config.id and acl.subject = any($1::int[]) and acl.access&$2 = $2)"), rbac.Groups, rbac.Access)
	}
	return base
}

func (c *Config) getFields() []string {
	return []string{
		"object_config.id",
		"object_config.enabled",
		"object_config.days_to_store",
		"object_config.period",
		"object_config.next_upload_on",
		"object_config.object_id",
		"wbt_class.name as object_name",
		"object_config.storage_id",
		"file_backend_profiles.name as storage_name",
		"object_config.domain_id",
		"object_config.created_at",
		"object_config.created_by",
		"object_config.updated_at",
		"object_config.updated_by",
	}
}
