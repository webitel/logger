package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/webitel/logger/internal/storage"
	"github.com/webitel/logger/internal/storage/postgres/utils"
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

func (c *Log) Select(ctx context.Context, opt *model.SearchOptions, filters any) ([]*model.Log, error) {
	var (
		query string
		args  []any
	)
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	base, err := storage.ApplyFiltersToBuilderBulk(c.GetQueryBaseFromSearchOptions(opt), logFieldsFilterMap, filters)
	if err != nil {
		return nil, err
	}
	switch req := base.(type) {
	case sq.SelectBuilder:
		query, args, _ = req.ToSql()
	default:
		return nil, errors.New("base of query is of wrong type")
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res, err := utils.ScanRows(rows, c.GetScanPlan)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Log) Insert(ctx context.Context, log *model.Log, domainId int) error {
	_, err := c.InsertBulk(ctx, []*model.Log{log}, domainId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Log) CheckRecordExist(ctx context.Context, objectName string, recordId int32) (bool, error) {
	db, err := c.storage.Database()
	if err != nil {
		return false, err
	}
	table, ok := recordTableMap[objectName]
	if !ok {
		return false, errors.New("object does not exist")
	}
	base := sq.Select("id").From(table.Path).Where(sq.Eq{"id": recordId}).PlaceholderFormat(sq.Dollar)
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected >= 0, nil
}

func (c *Log) InsertBulk(ctx context.Context, logs []*model.Log, domainId int) (int, error) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return 0, appErr
	}
	base := sq.Insert("logger.log").Columns("date", "action", "user_id", "user_ip", "new_state", "record_id", "config_id", "object_name").PlaceholderFormat(sq.Dollar)
	for _, log := range logs {
		base = base.Values(log.Date, log.Action, log.User.Id, log.UserIp, log.NewState, log.Record.Id, sq.Expr("(SELECT object_config.id FROM logger.object_config INNER JOIN directory.wbt_class ON object_config.object_id = wbt_class.id WHERE object_config.domain_id = ? AND wbt_class.name = ?)", domainId, log.Object.Name), log.Object.Name)
	}
	res, err := base.RunWith(db).ExecContext(ctx)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	return int(affected), err

}

func (c *Log) Delete(ctx context.Context, earlierThan time.Time, configId int) (int, error) {
	db, err := c.storage.Database()
	if err != nil {
		return 0, err
	}

	query := `DELETE FROM logger.log WHERE log.date < $1 AND log.config_id = $2 `
	rows, err := db.ExecContext(
		ctx,
		query,
		earlierThan,
		configId,
	)
	if err != nil {
		return 0, err
	}
	affected, err := rows.RowsAffected()
	return int(affected), err
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

func (c *Log) GetScanPlan(columns []string) []func(*model.Log) any {
	var binds []func(*model.Log) any
	for _, v := range columns {
		var bind func(*model.Log) any
		switch v {
		case "id":
			bind = func(dst *model.Log) any { return &dst.Id }
		case "date":
			bind = func(dst *model.Log) any { return &dst.Date }
		case "user_id":
			bind = func(dst *model.Log) any { return &dst.User.Id }
		case "user_name":
			bind = func(dst *model.Log) any { return &dst.User.Name }
		case "user_ip":
			bind = func(dst *model.Log) any { return &dst.UserIp }
		case "new_state":
			bind = func(dst *model.Log) any { return &dst.NewState }
		case "record_id":
			bind = func(dst *model.Log) any { return &dst.Record.Id }
		case "record_name":
			bind = func(dst *model.Log) any { return &dst.Record.Name }
		case "action":
			bind = func(dst *model.Log) any { return &dst.Action }
		case "config_id":
			bind = func(dst *model.Log) any { return &dst.ConfigId }
		case "object_id":
			bind = func(dst *model.Log) any { return &dst.Object.Id }
		case "object_name":
			bind = func(dst *model.Log) any { return &dst.Object.Name }
		default:
			bind = func(dst *model.Log) any { return nil }
		}
		binds = append(binds, bind)
	}
	return binds
}
