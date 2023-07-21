package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	errors "github.com/webitel/engine/model"
	"strings"
	"webitel_logger/model"
	"webitel_logger/storage"
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

func (c *Log) GetByObjectId(ctx context.Context, opt *model.SearchOptions, domainId int, objectId int) (*[]model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBaseFromSearchOptions(opt).Where(sq.Eq{"object_id": objectId}).Where(sq.Eq{"domain_id": domainId})
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.get_by_object_id.query_execute.fail", err.Error())
	}
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return &res, nil
}

func (c *Log) GetByUserId(ctx context.Context, opt *model.SearchOptions, userId int) (*[]model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBaseFromSearchOptions(opt).Where(
		sq.Eq{"userId": userId},
	)

	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.get_by_object_id.query_execute.fail", err.Error())
	}
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return &res, nil
}

func (c *Log) Insert(ctx context.Context, log *model.Log) (*model.Log, errors.AppError) {
	var newModel model.Log
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	rows, err := db.QueryContext(ctx,
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
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.insert.scan.error", err.Error())
	}
	res, appErr := c.ScanRows(rows)
	if err != nil {
		return nil, appErr
	}
	newModel = res[0]
	return &newModel, nil
}

func (c *Log) ScanRows(rows *sql.Rows) ([]model.Log, errors.AppError) {
	if rows == nil {
		return nil, errors.NewInternalError("postgres.log.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.scan.get_columns.error", err.Error())
	}
	var logs []model.Log

	for rows.Next() {
		var log model.Log
		binds := make([]func(dst *model.Log) interface{}, 0, 0)
		for _, v := range cols {

			switch v {
			case "id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Id })
			case "date":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Date })
			case "user_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.UserId })
			case "user_ip":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.UserIp })
			case "object_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.ObjectId })
			case "new_state":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.NewState })
			case "domain_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.DomainId })
			case "record_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.RecordId })
			case "action":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Action })
			default:
				panic("postgres.log.scan.get_columns.error: columns gotten from sql don't respond to model columns. Unknown column: " + v)
			}

		}
		bindFunc := func(binds []func(*model.Log) any) []any {
			var fields []any
			for _, v := range binds {
				fields = append(fields, v(&log))
			}
			return fields
		}
		err = rows.Scan(bindFunc(binds)...)
		if err != nil {
			return nil, errors.NewInternalError("postgres.log.scan.scan.error", err.Error())
		}
		logs = append(logs, log)
	}
	if len(logs) == 0 {
		return nil, errors.NewBadRequestError("postgres.log.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return logs, nil
}

func (c *Log) GetQueryBaseFromSearchOptions(opt *model.SearchOptions) sq.SelectBuilder {
	var fields []string
	if len(opt.Fields) == 0 {
		fields = append(fields,
			"id",
			"date",
			"user_id",
			"user_ip",
			"object_id",
			"new_state",
			"domain_id",
			"record_id",
			"action")
	} else {
		fields = opt.Fields
		if !contains(fields, "id") {
			fields = append(fields, "id")
		}
	}
	base := sq.Select(fields...).From("logger.log")
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
