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

func (c *Log) GetByObjectId(ctx context.Context, objectId int, domainId int) (*[]model.Log, errors.AppError) {
	var logs []model.Log
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	rows, err := db.QueryContext(ctx,
		`SELECT 
			id, 
			date, 
			user_id, 
			user_ip,
			object_id, 
			new_state, 
			domain_id,
			record_id,
		FROM 
			logger.log 
		WHERE 
			object_id = $1 
			AND domain_id = $2`,
		objectId, domainId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.get_by_object_id.query_execute.fail", err.Error())
	}
	for rows.Next() {
		var log model.Log
		err := rows.Scan(
			&log.Id,
			&log.Date,
			&log.UserId,
			&log.UserIp,
			&log.ObjectId,
			&log.NewState,
			&log.DomainId,
			&log.RecordId,
		)
		if err != nil {
			return nil, errors.NewInternalError("postgres.log.get_by_object_id.scan.fail", err.Error())
		}
		logs = append(logs, log)
	}
	return &logs, nil
}

func (c *Log) GetByUserId(ctx context.Context, userId int) (*[]model.Log, errors.AppError) {
	var results []model.Log
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	rows, err := db.QueryContext(ctx,
		`SELECT 
			id, 
			date, 
			action,
			user_id,
			user_ip,
			object_id, 
			new_state, 
			domain_id,
			record_id,
		FROM 
			logger.log 
		WHERE 
			user_id = $1`,
		userId)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.get_by_user_id.query_execution.fail", err.Error())
	}
	for rows.Next() {
		var log model.Log
		err := rows.Scan(&log.Id, &log.Date, &log.Action, &log.UserId, &log.UserIp, &log.ObjectId, &log.NewState, &log.DomainId, &log.RecordId)
		if err != nil {
			return nil, errors.NewInternalError("postgres.log.get_by_user_id.scan.fail", err.Error())
		}
		results = append(results, log)
	}
	return &results, nil
}

func (c *Log) Insert(ctx context.Context, log *model.Log) (*model.Log, errors.AppError) {
	var newModel model.Log
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	res := db.QueryRowContext(ctx,
		`INSERT INTO
			logger.log(date, action, user_id, user_ip, object_id, new_state, domain_id, record_id) 
		VALUES
			(
			$1, $2, $3, $4, $5, $6, $7, $8
			)
		RETURNING 
			id, date, action, user_id, user_ip, object_id, new_state, domain_id, record_id`,
		log.Date, log.Action, log.UserId, log.UserIp, log.ObjectId, log.NewState, log.DomainId, log.RecordId,
	)
	err := res.Scan(
		&newModel.Id,
		&newModel.Date,
		&newModel.Action,
		&newModel.UserId,
		&newModel.UserIp,
		&newModel.ObjectId,
		&newModel.NewState,
		&newModel.DomainId,
		&newModel.RecordId,
	)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.insert.scan.error", err.Error())
	}
	return &newModel, nil
}
