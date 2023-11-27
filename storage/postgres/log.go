package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/storage"
	"github.com/webitel/wlog"
)

var (
	logFieldsMap = map[string]string{
		"id":        "log.id",
		"date":      "log.date",
		"user":      "log.user_id, coalesce(wbt_user.name::varchar, wbt_user.username::varchar) as user_name",
		"user_ip":   "log.user_ip",
		"new_state": "log.new_state",
		"record":    "log.record_id",
		"action":    "log.action",
		"config_id": "log.config_id",
		"object":    "object_config.object_id, log.object_name",
	}
	recordTableMap = map[string]*storage.Table{
		"cc_queue": {
			Path:       "call_center.cc_queue",
			NameColumn: "name",
		},
		"scheme": {
			Path:       "flow.acr_routing_scheme",
			NameColumn: "name",
		},
	}
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

func (c *Log) Get(ctx context.Context, opt *model.SearchOptions, filters any) ([]*model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := ApplyFiltersToBuilder(c.GetQueryBaseFromSearchOptions(opt), filters)
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.get_by_object_id.query_execute.fail", err.Error())
	}
	defer rows.Close()
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (c *Log) CheckRecordExist(ctx context.Context, objectName string, recordId int32) (bool, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return false, appErr
	}
	table, ok := recordTableMap[objectName]
	if !ok {
		return false, errors.NewBadRequestError("postgres.log.check_record.invalid_args.error", "object does not exist")
	}
	base := sq.Select("id").From(table.Path).Where(sq.Eq{"id": recordId}).PlaceholderFormat(sq.Dollar)
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return false, errors.NewInternalError("postgres.log.check_record.query_execute.fail", err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.NewBadRequestError("postgres.log.check_record.get_res.fail", err.Error())
	}
	if rowsAffected <= 0 {
		return false, nil
	}
	return true, nil
}

func (c *Log) Insert(ctx context.Context, log *model.Log) errors.AppError {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return appErr
	}
	query := `INSERT INTO
			logger.log(date, action, user_id, user_ip, new_state, record_id, config_id, object_name)
			 $1, $2, $3, $4, $5, $6, $7, $8`
	wlog.Debug(query)
	_, err := db.ExecContext(ctx,
		query, log.Date, log.Action, log.User.Id, log.UserIp, log.NewState, log.Record.Id, log.ConfigId, log.Object.Name,
	)
	if err != nil {
		return errors.NewInternalError("postgres.log.insert.scan.error", err.Error())
	}
	return nil
}

func (c *Log) InsertMany(ctx context.Context, logs []*model.Log) errors.AppError {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return appErr
	}
	base := sq.Insert("logger.log").Columns("date", "action", "user_id", "user_ip", "new_state", "record_id", "config_id", "object_name").PlaceholderFormat(sq.Dollar)
	for _, log := range logs {
		base = base.Values(log.Date, log.Action, log.User.Id, log.UserIp, log.NewState, log.Record.Id, log.ConfigId, log.Object.Name)
	}
	sql, _, _ := base.ToSql()
	wlog.Debug(sql)
	_, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return errors.NewInternalError("postgres.log.insert.query.error", err.Error())
	}
	return nil
}

func (c *Log) DeleteByLowerThanDate(ctx context.Context, date time.Time, configId int) (int, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return 0, appErr
	}
	query := `DELETE FROM logger.log WHERE log.date < $1 AND log.config_id = $2 `
	wlog.Debug(query)
	rows, err := db.ExecContext(ctx, query, date, configId)
	if err != nil {
		return 0, errors.NewInternalError("postgres.log.delete_by_lowe_that_date.query.error", err.Error())
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return 0, errors.NewInternalError("postgres.log.delete_by_lowe_that_date.result.error", err.Error())
	}
	return int(affected), nil
}

func (c *Log) ScanRows(rows *sql.Rows) ([]*model.Log, errors.AppError) {
	if rows == nil {
		return nil, errors.NewInternalError("postgres.log.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.scan.get_columns.error", err.Error())
	}
	var logs []*model.Log

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
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.User.Id })
			case "user_name":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.User.Name })
			case "user_ip":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.UserIp })
			case "new_state":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.NewState })
			case "record_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Record.Id })
			case "action":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Action })
			case "config_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.ConfigId })
			case "object_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Object.Id })
			case "object_name":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Object.Name })
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
		logs = append(logs, &log)
	}
	if len(logs) == 0 {
		return nil, errors.NewBadRequestError("postgres.log.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return logs, nil
}

func (c *Log) GetQueryBaseFromSearchOptions(opt *model.SearchOptions) sq.SelectBuilder {
	var fields []string
	if opt == nil {
		return c.GetQueryBase(c.getFields())
	}
	for _, v := range opt.Fields {
		fields = append(fields, logFieldsMap[v])
	}
	if len(fields) == 0 {
		fields = append(fields,
			c.getFields()...)
	}
	base := c.GetQueryBase(fields)
	if opt.Search != "" {
		base = base.Where(sq.Like{"user_ip": opt.Search + "%"})
	}
	if opt.Sort != "" {
		splitted := strings.Split(opt.Sort, ":")
		if len(splitted) == 2 {
			order := splitted[0]
			column := splitted[1]
			if column == "user" {
				column = "user_name"
			}
			base = base.OrderBy(fmt.Sprintf("%s %s", column, order))
		}

	}
	offset := (opt.Page - 1) * opt.Size
	if offset < 0 {
		offset = 0
	}
	return base.Offset(uint64(offset)).Limit(uint64(opt.Size + 1))
}

func (c *Log) GetQueryBase(fields []string) sq.SelectBuilder {
	base := sq.Select(fields...).
		From("logger.log").
		JoinClause("LEFT JOIN directory.wbt_user ON wbt_user.id = log.user_id").
		JoinClause("LEFT JOIN logger.object_config ON object_config.id = log.config_id").
		PlaceholderFormat(sq.Dollar)

	return base
}

func (c *Log) getFields() []string {
	var fields []string
	for _, value := range logFieldsMap {
		fields = append(fields, value)
	}
	return fields
}
