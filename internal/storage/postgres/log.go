package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/webitel/logger/internal/storage"
	"log/slog"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/logger/internal/model"
)

var (
	logFieldsSelectMap = map[string]string{
		model.LogFields.Id:       "log.id",
		model.LogFields.UserIp:   "log.user_ip",
		model.LogFields.Action:   "log.action",
		model.LogFields.ConfigId: "log.config_id",
		model.LogFields.Date:     "log.date",
		model.LogFields.NewState: "log.new_state",

		// [combined alias]
		model.LogFields.Object: "log.object_name, object_config.object_id",
		model.LogFields.User:   "log.user_id, coalesce(wbt_user.name::varchar, wbt_user.username::varchar) as user_name",
		model.LogFields.Record: "log.record_id",
	}
	logFieldsFilterMap = map[string]string{
		model.LogFields.Id:       "log.id",
		model.LogFields.User:     "log.user_id",
		model.LogFields.UserIp:   "log.user_ip",
		model.LogFields.Action:   "log.action",
		model.LogFields.ConfigId: "log.config_id",
		model.LogFields.Date:     "log.date",
		model.LogFields.NewState: "log.new_state",
		model.LogFields.Object:   "object_config.object_id",
		model.LogFields.User:     "log.user_id",
		model.LogFields.Record:   "log.record_id",
	}

	recordTableMap = map[string]*storage.Table{
		"cc_queue": {
			Path:       "call_center.cc_queue",
			NameColumn: "name",
		},
		"schema": {
			Path:       "flow.acr_routing_scheme",
			NameColumn: "name",
		},
		"users": {
			Path:       "directory.wbt_user",
			NameColumn: "name",
		},
		"calendars": {
			Path:       "flow.calendar",
			NameColumn: "name",
		},
		"cc_list": {
			Path:       "call_center.cc_list",
			NameColumn: "name",
		},
		"cc_team": {
			Path:       "call_center.cc_team",
			NameColumn: "name",
		},
		// no name field??
		//"cc_agent": {
		//	Path:       "call_center.cc_agent",
		//	NameColumn: "name",
		//},
		"cc_resource": {
			Path:       "call_center.cc_outbound_resource",
			NameColumn: "name",
		},
		"cc_resource_group": {
			Path:       "call_center.cc_outbound_resource_group",
			NameColumn: "name",
		},
		"chat_bots": {
			Path:       "chat.bot",
			NameColumn: "name",
		},
	}
)

func init() {
	v, ok := logFieldsSelectMap[model.LogFields.Record]
	if ok {
		v += ", (case "
		for objectName, table := range recordTableMap {
			v += fmt.Sprintf("when log.object_name = '%s' then (select %s.%s from %[2]s where id = record_id) ", objectName, table.Path, table.NameColumn)
		}
		v += " end) record_name"
	}
	logFieldsSelectMap[model.LogFields.Record] = v
}

type Log struct {
	storage storage.Storage
}

func newLogStore(store storage.Storage) (storage.LogStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.log.new_log.check.bad_arguments", "error creating log interface to the log table, main store is nil")
	}
	return &Log{storage: store}, nil
}

func (c *Log) Select(ctx context.Context, opt *model.SearchOptions, filters any) ([]*model.Log, model.AppError) {
	var (
		query string
		args  []any
	)
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base, appErr := storage.ApplyFiltersToBuilderBulk(c.GetQueryBaseFromSearchOptions(opt), logFieldsFilterMap, filters)
	if appErr != nil {
		return nil, appErr
	}
	switch req := base.(type) {
	case sq.SelectBuilder:
		query, args, _ = req.ToSql()
	default:
		return nil, model.NewInternalError("store.sql_scheme_variable.get.base_type.wrong", "base of query is of wrong type")
	}
	rows, err := db.QueryContext(ctx, query, args...)
	slog.Debug(query, args)
	if err != nil {
		return nil, model.NewInternalError("postgres.log.get_by_object_id.query_execute.fail", err.Error())
	}
	defer rows.Close()
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (c *Log) Insert(ctx context.Context, log *model.Log, domainId int) model.AppError {
	err := c.InsertBulk(ctx, []*model.Log{log}, domainId)
	if err != nil {
		return model.NewInternalError("postgres.log.insert.scan.error", err.Error())
	}
	return nil
}

func (c *Log) CheckRecordExist(ctx context.Context, objectName string, recordId int32) (bool, model.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return false, appErr
	}
	table, ok := recordTableMap[objectName]
	if !ok {
		return false, model.NewBadRequestError("postgres.log.check_record.invalid_args.error", "object does not exist")
	}
	base := sq.Select("id").From(table.Path).Where(sq.Eq{"id": recordId}).PlaceholderFormat(sq.Dollar)
	res, err := base.RunWith(db).ExecContext(ctx)
	slog.Debug(base.ToSql())
	if err != nil {
		return false, model.NewInternalError("postgres.log.check_record.query_execute.fail", err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, model.NewBadRequestError("postgres.log.check_record.get_res.fail", err.Error())
	}
	if rowsAffected <= 0 {
		return false, nil
	}
	return true, nil
}

func (c *Log) InsertBulk(ctx context.Context, logs []*model.Log, domainId int) model.AppError {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return appErr
	}
	base := sq.Insert("logger.log").Columns("date", "action", "user_id", "user_ip", "new_state", "record_id", "config_id", "object_name").PlaceholderFormat(sq.Dollar)
	for _, log := range logs {
		base = base.Values(log.Date, log.Action, log.User.Id, log.UserIp, log.NewState, log.Record.Id, sq.Expr("(SELECT object_config.id FROM logger.object_config INNER JOIN directory.wbt_class ON object_config.object_id = wbt_class.id WHERE object_config.domain_id = ? AND wbt_class.name = ?)", domainId, log.Object.Name), log.Object.Name)
	}

	_, err := base.RunWith(db).ExecContext(ctx)
	slog.Debug(base.ToSql())
	if err != nil {
		return model.NewInternalError("postgres.log.insert.query.error", err.Error())
	}
	return nil
}

func (c *Log) Delete(ctx context.Context, date time.Time, configId int) (int, model.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return 0, appErr
	}

	query := `DELETE FROM logger.log WHERE log.date < $1 AND log.config_id = $2 `
	slog.Debug(query)
	rows, err := db.ExecContext(
		ctx,
		query,
		date,
		configId,
	)
	if err != nil {
		return 0, model.NewInternalError("postgres.log.delete_by_lowe_that_date.query.error", err.Error())
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return 0, model.NewInternalError("postgres.log.delete_by_lowe_that_date.result.error", err.Error())
	}
	return int(affected), nil
}

func (c *Log) ScanRows(rows *sql.Rows) ([]*model.Log, model.AppError) {
	if rows == nil {
		return nil, model.NewInternalError("postgres.log.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, model.NewInternalError("postgres.log.scan.get_columns.error", err.Error())
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
			case "record_name":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Record.Name })
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
			return nil, model.NewInternalError("postgres.log.scan.scan.error", err.Error())
		}
		logs = append(logs, &log)
	}
	if len(logs) == 0 {
		return nil, model.NewBadRequestError("postgres.log.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return logs, nil
}

func (c *Log) GetQueryBaseFromSearchOptions(opt *model.SearchOptions) sq.SelectBuilder {
	var fields []string
	if opt == nil {
		return c.GetQueryBase(c.getFields())
	}
	for _, v := range opt.Fields {
		if columnName, ok := logFieldsSelectMap[v]; ok {
			fields = append(fields, columnName)
		} else {
			fields = append(fields, v)
		}
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
			// Lookup columns -- order by name
			switch column {
			case model.LogFields.User:
				column = "user_name"
			case model.LogFields.Object:
				column = "object_name"
			case model.LogFields.Record:
				column = "record_name"
			}
			base = base.OrderBy(fmt.Sprintf("%s %s", column, order))
		}

	}
	offset := (opt.Page - 1) * opt.Size
	if offset < 0 {
		offset = 0
	}
	if opt.Size != 0 {
		base = base.Limit(uint64(opt.Size + 1))
	}
	return base.Offset(uint64(offset))
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
	for _, value := range logFieldsSelectMap {
		fields = append(fields, value)
	}
	return fields
}
