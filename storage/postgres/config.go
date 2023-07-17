package postgres

import (
	"context"
	"database/sql"
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
			storage_id = $5,
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
	//var nextUploadOn time.Time
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

func (c *Config) GetByObjectId(ctx context.Context, objId int, domainId int) (*model.Config, errors.AppError) {
	var conf model.Config
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	res := db.QueryRow(
		`SELECT
			id,
			enabled,
			days_to_store,
			period,
			next_upload_on,
			object_id,
			storage_id,
			domain_id 
		FROM
			logger.object_config 
		WHERE
			object_id = $1 
			AND domain_id = $2`,
		objId, domainId)
	err := res.Scan(&conf.Id, &conf.Enabled, &conf.DaysToStore, &conf.Period, &conf.NextUploadOn, &conf.ObjectId, &conf.StorageId, &conf.DomainId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_object_id.scan.fail", err.Error())
	}
	return &conf, nil
}

func (c *Config) GetById(ctx context.Context, id int) (*model.Config, errors.AppError) {
	var conf model.Config
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	res := db.QueryRow(
		`SELECT
			id,
			enabled,
			days_to_store,
			period,
			next_upload_on,
			object_id,
			storage_id,
			domain_id 
		FROM
			logger.object_config 
		WHERE
			id = $1`,
		id)
	err := res.Scan(&conf.Id, &conf.Enabled, &conf.DaysToStore, &conf.Period, &conf.NextUploadOn, &conf.ObjectId, &conf.StorageId, &conf.DomainId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_by_object_id.scan.fail", err.Error())
	}
	return &conf, nil
}

func (c *Config) GetAll(ctx context.Context, domainId int) (*[]model.Config, errors.AppError) {
	var configs []model.Config
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	rows, err := db.Query(
		`SELECT
			id,
			enabled,
			days_to_store,
			period,
			next_upload_on,
			object_id,
			storage_id,
			domain_id 
		FROM
			logger.object_config 
		WHERE
			domain_id = $1;`,
		domainId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.get_all.query.fail", err.Error())
	}
	for rows.Next() {
		var config model.Config
		err := rows.Scan(&config.Id, &config.Enabled, &config.DaysToStore, &config.Period, &config.NextUploadOn, &config.ObjectId, &config.StorageId, &config.DomainId)
		if err != nil {
			return nil, errors.NewInternalError("postgres.config.get_all.scan.fail", err.Error())
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil, errors.NewInternalError("postgres.config.get_all.result_check.empty", sql.ErrNoRows.Error())
	}
	return &configs, nil
}
