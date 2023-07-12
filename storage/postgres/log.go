package postgres

import (
	"context"
	"webitel_logger/model"
	"webitel_logger/storage"

	errors "github.com/webitel/engine/model"
)

type Log struct {
	storage storage.Storage
}

func newLogStore(store storage.Storage) (storage.LogStore, errors.AppError) {
	if store == nil {
		return nil, errors.NewInternalError("postgres.log.new_log.check.bad_arguments", "error creating log interface to the log table, main store is nil")
	}
	return &Log{storage: store}, nil
}

func (c *Log) GetByObjectId(ctx context.Context, objectId int) (*[]model.Log, errors.AppError) {
	_, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *Log) GetByUserId(ctx context.Context, userId int) (*[]model.Log, errors.AppError) {
	_, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *Log) Insert(ctx context.Context, log *model.Log) (*model.Log, errors.AppError) {
	_, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
