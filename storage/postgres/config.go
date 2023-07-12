package postgres

import (
	"context"
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
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	//var nextUploadOn time.Time
	_, err := db.ExecContext(ctx, `UPDATE config 
	SET enabled = :Enabled,
	days_to_store = :DaysToStore,
	period = :Period,
	next_upload_on = :NextUploadOn,
	storage_id = :StorageId
	WHERE id = :Id;`, map[string]any{
		":Enabled":      conf.Enabled,
		":DaysToStore":  conf.DaysToStore,
		":Period":       conf.Period,
		":NextUploadOn": conf.NextUploadOn,
		":StorageId":    conf.StorageId,
	})
	if err != nil {
		return nil, errors.NewInternalError("postgres.config.update.query_execution.fail", err.Error())
	}
	return nil, nil
}

func (c *Config) GetByObjectId(ctx context.Context, objId int) (*model.Config, errors.AppError) {
	_, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *Config) GetAll(ctx context.Context) (*[]model.Config, errors.AppError) {
	_, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
