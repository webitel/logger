package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	errors "github.com/webitel/engine/model"
	"strings"
	"time"
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
	base := c.GetQueryBaseFromSearchOptions(opt).Where(
		sq.Eq{"object_config.domain_id": domainId},
		sq.Eq{"object_config.object_id": objectId},
	).JoinClause("LEFT JOIN logger.object_config ON log.config_id = object_config.id")
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

func (c *Log) GetByObjectIdWithDates(ctx context.Context, domainId int, objectId int, dateFrom time.Time, dateTo time.Time) (*[]model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBase(c.getFields()).Where(
		sq.Eq{"object_config.domain_id": domainId},
		sq.Eq{"object_config.object_id": objectId},
		sq.GtOrEq{"log.date": dateFrom},
		sq.LtOrEq{"log.date": dateTo},
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

func (c *Log) GetByConfigId(ctx context.Context, opt *model.SearchOptions, configId int) (*[]model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBaseFromSearchOptions(opt).Where(sq.Eq{"log.config_id": configId})
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

func (c *Log) GetByConfigIdWithDates(ctx context.Context, configId int, dateFrom time.Time, dateTo time.Time) (*[]model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBase(c.getFields()).Where(
		sq.Eq{"log.config_id": configId},
		sq.GtOrEq{"log.date": dateFrom},
		sq.LtOrEq{"log.date": dateTo},
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

func (c *Log) GetByUserId(ctx context.Context, opt *model.SearchOptions, userId int) (*[]model.Log, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := c.GetQueryBaseFromSearchOptions(opt).Where(
		sq.Eq{"user_id": userId},
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
			logger.log(date, action, user_id, user_ip, new_state, record_id, config_id) 
		VALUES
			(
			$1, $2, $3, $4, $5, $6, $7
			)
		RETURNING 
			id, date, action, user_id, user_ip,  new_state, record_id, config_id`,
		log.Date, log.Action, log.User.Id, log.UserIp, log.NewState, log.RecordId, log.ConfigId,
	)
	if err != nil {
		return nil, errors.NewInternalError("postgres.log.insert.scan.error", err.Error())
	}
	res, appErr := c.ScanRows(rows)
	if appErr != nil {
		return nil, appErr
	}
	newModel = res[0]
	return &newModel, nil
}

func (c *Log) DeleteByLowerThanDate(ctx context.Context, date time.Time, configId int) (int, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return 0, appErr
	}
	rows, err := db.ExecContext(ctx,
		`DELETE FROM logger.log WHERE log.date < $1 AND log.config_id = $2 `,
		date, configId,
	)
	if err != nil {
		return 0, errors.NewInternalError("postgres.log.delete_by_lowe_that_date.query.error", err.Error())
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return 0, errors.NewInternalError("postgres.log.delete_by_lowe_that_date.result.error", err.Error())
	}
	return int(affected), nil
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
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.User.Id })
			case "user_name":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.User.Name })
			case "user_ip":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.UserIp })
			case "new_state":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.NewState })
			case "record_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.RecordId })
			case "action":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.Action })
			case "config_id":
				binds = append(binds, func(dst *model.Log) interface{} { return &dst.ConfigId })
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
	for _, v := range opt.Fields {
		switch v {
		case "id":
			fields = append(fields, "log.id")
		case "date":
			fields = append(fields, "log.date")
		case "user_ip":
			fields = append(fields, "log.user_ip")
		case "new_state":
			fields = append(fields, "log.new_state")
		case "record_id":
			fields = append(fields, "log.record_id")
		case "action":
			fields = append(fields, "log.action")
		case "config_id":
			fields = append(fields, "log.config_id")
		case "user":
			fields = append(fields, "coalesce(wbt_user.name::varchar, wbt_user.username::varchar) as user_name")
			fields = append(fields, "log.user_id")
		}
	}
	if len(fields) == 0 {
		fields = append(fields,
			c.getFields()...)
	}
	base := c.GetQueryBase(fields)
	//if opt.Search != "" {
	//	base = base.Where(sq.Like{"description": "%" + strings.ToLower(opt.Search) + "%"})
	//}
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
	return base.Offset(uint64(offset)).Limit(uint64(opt.Size)).PlaceholderFormat(sq.Dollar)
}

func (c *Log) GetQueryBase(fields []string) sq.SelectBuilder {
	return sq.Select(fields...).From("logger.log").JoinClause("LEFT JOIN directory.wbt_user ON wbt_user.id = log.user_id")
}

func (c *Log) getFields() []string {
	return []string{
		"log.id",
		"log.date",
		"log.user_id",
		"coalesce(wbt_user.name::varchar, wbt_user.username::varchar) as user_name",
		"log.user_ip",
		"log.new_state",
		"log.record_id",
		"log.action",
		"log.config_id",
	}
}
