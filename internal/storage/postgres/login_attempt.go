package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/webitel/logger/internal/model"
	"github.com/webitel/logger/internal/storage"
	"github.com/webitel/logger/internal/storage/postgres/utils"
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
func newLoginAttemptStore(store storage.Storage) (storage.LoginAttemptStore, error) {
	if store == nil {
		return nil, errors.New("error creating login attempt store interface, main store is nil")
	}
	return &LoginAttemptStore{storage: store}, nil
}

// endregion

func (l *LoginAttemptStore) Insert(ctx context.Context, attempt *model.LoginAttempt) (*model.LoginAttempt, error) {
	base := squirrel.Insert(loginAttemptTable).Columns("date", "domain_id", "domain_name", "success", "auth_type", "user_id", "user_name", "user_ip", "user_agent", "details").
		Values(attempt.Date, squirrel.Expr("coalesce(?, (select dc from directory.wbt_domain where name = ? limit 1))", attempt.DomainId, attempt.DomainName), attempt.DomainName, attempt.Success, attempt.AuthType, squirrel.Expr("coalesce(?, (select id from directory.wbt_user where username = ? limit 1))", attempt.UserId, attempt.UserName), attempt.UserName, attempt.UserIp, attempt.UserAgent, attempt.Details).Suffix("returning *").PlaceholderFormat(squirrel.Dollar)
	db, err := l.storage.Database()
	if err != nil {
		return nil, err
	}
	sql, params, _ := base.ToSql()
	rows, err := db.QueryContext(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	response, err := utils.ScanRow(rows, l.GetScanPlan)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func (l *LoginAttemptStore) Select(ctx context.Context, searchOpts *model.SearchOptions, filters any) ([]*model.LoginAttempt, error) {
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
		return nil, errors.New("invalid search options")
	}

	rows, err := db.QueryContext(ctx, query, params...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	res, err := utils.ScanRows(rows, l.GetScanPlan)
	if err != nil {
		return nil, err
	}
	return res, nil
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

func (c *LoginAttemptStore) GetScanPlan(columns []string) []func(*model.LoginAttempt) any {
	var binds []func(*model.LoginAttempt) any
	for _, v := range columns {
		var bind func(*model.LoginAttempt) any
		switch v {
		case "id":
			bind = func(dst *model.LoginAttempt) any { return &dst.Id }
		case "date":
			bind = func(dst *model.LoginAttempt) any { return &dst.Date }
		case "domain_id":
			bind = func(dst *model.LoginAttempt) any { return &dst.DomainId }
		case "domain_name":
			bind = func(dst *model.LoginAttempt) any { return &dst.DomainName }
		case "success":
			bind = func(dst *model.LoginAttempt) any { return &dst.Success }
		case "auth_type":
			bind = func(dst *model.LoginAttempt) any { return &dst.AuthType }
		case "user_id":
			bind = func(dst *model.LoginAttempt) any { return &dst.UserId }
		case "user_name":
			bind = func(dst *model.LoginAttempt) any { return &dst.UserName }
		case "user_ip":
			bind = func(dst *model.LoginAttempt) any { return &dst.UserIp }
		case "user_agent":
			bind = func(dst *model.LoginAttempt) any { return &dst.UserAgent }
		case "details":
			bind = func(dst *model.LoginAttempt) any { return &dst.Details }
		default:
			continue
		}
		binds = append(binds, bind)
	}
	return binds
}
