package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strings"
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

func (c *Config) Update(ctx context.Context, conf *model.Config) (*model.Config, errors.AppError) {
	var newConfig model.Config
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	row := db.QueryRowContext(ctx,
		`UPDATE 
			logger.object_config 
		SET 
			enabled = $1, 
			days_to_store = $2, 
			period = $3, 
			next_upload_on = $4, 
			storage_id = $5
		WHERE 
			object_id = $6 AND domain_id = $7
		RETURNING id, enabled, days_to_store, period, next_upload_on, storage_id, domain_id, object_id;`,
		conf.Enabled, conf.DaysToStore, conf.Period, conf.NextUploadOn, conf.StorageId, conf.ObjectId, conf.DomainId)
	err := row.Scan(&newConfig.Id, &newConfig.Enabled, &newConfig.DaysToStore, &newConfig.Period, &newConfig.NextUploadOn, &newConfig.StorageId, &newConfig.DomainId, &newConfig.ObjectId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.update.query_execution.fail", err.Error())
	}
	return &newConfig, nil
}

func (c *Config) Insert(ctx context.Context, conf *model.Config) (*model.Config, errors.AppError) {
	var newConfig model.Config
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	res := db.QueryRowContext(ctx,
		`INSERT INTO
		logger.object_config(enabled, days_to_store, period, next_upload_on, storage_id, domain_id, object_id) 
		VALUES
			(
			$1, $2, $3, $4, $5, $6, $7
			)
		RETURNING id, enabled, days_to_store, period, next_upload_on, storage_id, domain_id, object_id`,
		conf.Enabled,
		conf.DaysToStore,
		conf.Period,
		conf.NextUploadOn,
		conf.StorageId,
		conf.DomainId,
		conf.ObjectId)
	err := res.Scan(&newConfig.Id, &newConfig.Enabled, &newConfig.DaysToStore, &newConfig.Period, &newConfig.NextUploadOn, &newConfig.StorageId, &newConfig.DomainId, &newConfig.ObjectId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.insert.scan_new.fail", err.Error())
	}
	return &newConfig, nil
}

func (c *Config) GetByObjectId(ctx context.Context /*opt *model.SearchOptions, */, domainId int, objId int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	//base := c.GetQueryBaseFromSearchOptions(opt)
	//rows, err := base.Where(sq.Eq{"domain_id": domainId}).Where(sq.Eq{"object_id": objId}).RunWith(db).QueryContext(ctx)
	rows, err := sq.Select("id",
		"enabled",
		"days_to_store",
		"period",
		"next_upload_on",
		"object_id",
		"storage_id",
		"domain_id").From("logger.object_config").Where(sq.Eq{"domain_id": domainId}).Where(sq.Eq{"object_id": objId}).PlaceholderFormat(sq.Dollar).RunWith(db).QueryContext(ctx)

	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_object.query.fail", err.Error())
	}
	configs, appErr := c.ScanRow(rows)
	if appErr != nil {
		return nil, appErr
	}
	return configs, nil
}

func (c *Config) GetById(ctx context.Context, opt *model.SearchOptions, id int) (*model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBaseFromSearchOptions(opt)
	rows, err := base.Where(sq.Eq{"id": id}).RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_id.query.fail", err.Error())
	}
	config, appErr := c.ScanRow(rows)
	if appErr != nil {
		return nil, appErr
	}
	return config, nil
}

func (c *Config) GetAll(ctx context.Context, opt *model.SearchOptions, domainId int) (*[]model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBaseFromSearchOptions(opt)

	rows, err := base.Where(sq.Eq{"domain_id": domainId}).RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_all.query.fail", err.Error())
	}
	configs, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return &configs, nil
}

func (c *Config) GetAllEnabledConfigs(ctx context.Context) (*[]model.Config, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	rows, err := sq.Select("object_config.id",
		"object_config.enabled",
		"object_config.days_to_store",
		"object_config.period",
		"object_config.next_upload_on",
		"object_config.object_id",
		"object_config.storage_id",
		"object_config.domain_id").From("logger.object_config").Where(sq.Eq{"enabled": true}).PlaceholderFormat(sq.Dollar).RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_all_enabled_configs.query.fail", err.Error())
	}
	configs, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return &configs, nil
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
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.ObjectId })
			case "storage_id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.StorageId })
			case "domain_id":
				binds = append(binds, func(dst *model.Config) interface{} { return &dst.DomainId })
			default:
				panic("postgres.log.scan.get_columns.error: columns gotten from sql don't respond to model columns. Unknown column: " + v)
			}

		}
		bindFunc := func(binds []func(config2 *model.Config) any) []any {
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

func (c *Config) GetQueryBaseFromSearchOptions(opt *model.SearchOptions) sq.SelectBuilder {
	var fields []string
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
		case "object_id":
			fields = append(fields, "object_config.object_id")
		case "storage_id":
			fields = append(fields, "object_config.storage_id")
		case "domain_id":
			fields = append(fields, "object_config.domain_id")
		}
	}
	if len(fields) == 0 {
		fields = append(fields,
			"object_config.id",
			"object_config.enabled",
			"object_config.days_to_store",
			"object_config.period",
			"object_config.next_upload_on",
			"object_config.object_id",
			"object_config.storage_id",
			"object_config.domain_id")
	}
	base := sq.Select(fields...).From("logger.object_config")
	//if opt.Search != "" {
	//	base = base.Where(sq.Like{"description": "%" + strings.ToLower(opt.Search) + "%"})
	//}
	if opt.Sort != "" {
		splitted := strings.Split(opt.Sort, ":")
		if len(splitted) == 2 {
			order := splitted[0]
			column := splitted[1]
			base = base.OrderBy(fmt.Sprintf("%s %s", column, order))
		}

	}
	offset := (opt.Page - 1) * opt.Size
	if offset < 0 {
		offset = 0
	}
	return base.Offset(uint64(offset)).Limit(uint64(opt.Size)).PlaceholderFormat(sq.Dollar)
}
