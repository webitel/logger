package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/storage"
	"strings"
)

var (
	loginAttemptSelectMap = map[string]string{
		model.LoginAttemptFields.Id:         "login_attempt.id",
		model.LoginAttemptFields.Date:       "login_attempt.date",
		model.LoginAttemptFields.DomainId:   "login_attempt.domain_id",
		model.LoginAttemptFields.Success:    "login_attempt.success",
		model.LoginAttemptFields.AuthType:   "login_attempt.auth_type",
		model.LoginAttemptFields.UserId:     "login_attempt.user_id",
		model.LoginAttemptFields.UserName:   "coalesce(wbt_user.name, wbt_user.username)",
		model.LoginAttemptFields.UserIp:     "login_attempt.user_ip",
		model.LoginAttemptFields.UserAgent:  "login_attempt.user_agent",
		model.LoginAttemptFields.DomainName: "login_attempt.domain_name",
		model.LoginAttemptFields.Details:    "login_attempt.details",

		// alias
		model.LoginAttemptFields.User: "login_attempt.user_id, coalesce(wbt_user.name, wbt_user.username)",
	}
	loginAttemptFilterMap = map[string]string{
		model.LoginAttemptFields.Id:         "login_attempt.id",
		model.LoginAttemptFields.Date:       "login_attempt.date",
		model.LoginAttemptFields.DomainId:   "login_attempt.domain_id",
		model.LoginAttemptFields.Success:    "login_attempt.success",
		model.LoginAttemptFields.AuthType:   "login_attempt.auth_type",
		model.LoginAttemptFields.User:       "wbt_user.name",
		model.LoginAttemptFields.UserName:   "coalesce(wbt_user.name, wbt_user.username)",
		model.LoginAttemptFields.UserId:     "login_attempt.user_id",
		model.LoginAttemptFields.UserIp:     "login_attempt.user_ip",
		model.LoginAttemptFields.UserAgent:  "login_attempt.user_agent",
		model.LoginAttemptFields.DomainName: "login_attempt.domain_name",
		model.LoginAttemptFields.Details:    "login_attempt.details",
	}
)

type LoginAttemptStore struct {
	storage storage.Storage
}

const (
	loginAttemptTable = "logger.login_attempt"
)

// region CONSTRUCTOR
func newLoginAttemptStore(store storage.Storage) (storage.LoginAttemptStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.login_attempt.new_config.check.bad_arguments", "error creating login attempt store interface, main store is nil")
	}
	return &LoginAttemptStore{storage: store}, nil
}

// endregion

func (l *LoginAttemptStore) Insert(ctx context.Context, attempt *model.LoginAttempt) (*model.LoginAttempt, model.AppError) {
	base := squirrel.Insert(loginAttemptTable).Columns("date", "domain_id", "domain_name", "success", "auth_type", "user_id", "user_name", "user_ip", "user_agent", "details").
		Values(attempt.Date, squirrel.Expr("coalesce(?, (select dc from directory.wbt_domain where name = ? limit 1))", attempt.DomainId, attempt.DomainName), attempt.DomainName, attempt.Success, attempt.AuthType, squirrel.Expr("coalesce(?, (select id from directory.wbt_user where username = ? limit 1))", attempt.UserId, attempt.UserName), attempt.UserName, attempt.UserIp, attempt.UserAgent, attempt.Details).Suffix("returning *").PlaceholderFormat(squirrel.Dollar)

	db, appErr := l.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	sql, params, _ := base.ToSql()
	res, err := db.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, model.NewInternalError("postgres.login_attempt.insert.execute.error", err.Error())
	}
	response, appErr := l.ScanRows(res)
	if appErr != nil {
		return nil, appErr
	}
	return response[0], nil

}

func (l *LoginAttemptStore) Get(ctx context.Context, searchOpts *model.SearchOptions, filters any) ([]*model.LoginAttempt, model.AppError) {
	var (
		query  string
		params []any
	)
	db, err := l.storage.Database()
	if err != nil {
		return nil, err
	}
	base, err := storage.ApplyFiltersToBuilderBulk(l.GetQueryBaseFromSearchOptions(searchOpts), loginAttemptFilterMap, filters)
	if err != nil {
		return nil, err
	}
	switch req := base.(type) {
	case squirrel.SelectBuilder:
		query, params, _ = req.ToSql()
	default:
		return nil, model.NewInternalError("store.sql_scheme_variable.get.base_type.wrong", "base of query is of wrong type")
	}

	rows, defErr := db.QueryContext(ctx, query, params...)
	if defErr != nil {
		return nil, model.NewInternalError("postgres.login_attempt.select.execute_query.error", defErr.Error())
	}
	return l.ScanRows(rows)
}

func (c *LoginAttemptStore) ScanRows(rows *sql.Rows) ([]*model.LoginAttempt, model.AppError) {
	if rows == nil {
		return nil, model.NewInternalError("postgres.login_attempt.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, model.NewInternalError("postgres.login_attempt.scan.get_columns.error", err.Error())
	}
	var logs []*model.LoginAttempt

	for rows.Next() {
		var log model.LoginAttempt
		binds := make([]func(dst *model.LoginAttempt) interface{}, 0, 0)
		for _, v := range cols {

			switch v {
			case "id":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.Id })
			case "date":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.Date })
			case "domain_id":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.DomainId })
			case "domain_name":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.DomainName })
			case "success":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.Success })
			case "auth_type":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.AuthType })
			case "user_id":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.UserId })
			case "user_name":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.UserName })
			case "user_ip":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.UserIp })
			case "user_agent":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.UserAgent })
			case "details":
				binds = append(binds, func(dst *model.LoginAttempt) interface{} { return &dst.Details })
			default:
				panic("postgres.login_attempt.scan.get_columns.error: columns gotten from sql don't respond to model columns. Unknown column: " + v)
			}

		}
		bindFunc := func(binds []func(*model.LoginAttempt) any) []any {
			var fields []any
			for _, v := range binds {
				fields = append(fields, v(&log))
			}
			return fields
		}
		err = rows.Scan(bindFunc(binds)...)
		if err != nil {
			return nil, model.NewInternalError("postgres.login_attempt.scan.scan.error", err.Error())
		}
		logs = append(logs, &log)
	}
	if len(logs) == 0 {
		return nil, model.NewBadRequestError("postgres.login_attempt.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return logs, nil
}

func (c *LoginAttemptStore) GetQueryBaseFromSearchOptions(opt *model.SearchOptions) squirrel.SelectBuilder {
	var fields []string
	if opt == nil {
		return c.GetQueryBase(c.getFields())
	}
	for _, v := range opt.Fields {
		if columnName, ok := loginAttemptSelectMap[v]; ok {
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
		base = base.Where(squirrel.Like{loginAttemptFilterMap[model.LoginAttemptFields.UserName]: opt.Search + "%"})
	}
	if opt.Sort != "" {
		splitted := strings.Split(opt.Sort, ":")
		if len(splitted) == 2 {
			order := splitted[0]
			column := splitted[1]
			// Lookup columns -- order by name
			switch column {
			case model.LoginAttemptFields.UserName:
				column = "user_name"
			case model.LoginAttemptFields.Date:
				column = "date"
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

func (c *LoginAttemptStore) GetQueryBase(fields []string) squirrel.SelectBuilder {
	base := squirrel.Select(fields...).
		From(loginAttemptTable).
		JoinClause("LEFT JOIN directory.wbt_user ON wbt_user.id = login_attempt.user_id").
		PlaceholderFormat(squirrel.Dollar)

	return base
}

func (c *LoginAttemptStore) getFields() []string {
	var fields []string
	for _, value := range logFieldsSelectMap {
		fields = append(fields, value)
	}
	return fields
}
