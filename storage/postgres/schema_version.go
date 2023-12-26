package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/storage"
	"github.com/webitel/wlog"
	"strings"
)

var (
	schemaVersionFieldsMap = map[string]string{
		"id":          "schema_version.id",
		"schema":      "schema_version.schema_id",
		"created_on":  "schema_version.created_on",
		"created_by":  "schema_version.created_by as created_by_id, wbt_user.name as created_by_name",
		"object_data": "schema_version.object_data",
		"version":     "schema_version.version",
		"note":        "schema_version.note",
	}
)

const maxVersionsCountPerSchema = 10

type SchemaVersion struct {
	storage          storage.Storage
	maxVersionsCount int
}

func newSchemaVersionStore(store storage.Storage) (storage.SchemaVersionStore, errors.AppError) {
	if store == nil {
		return nil, errors.NewInternalError("postgres.log.new_log.check.bad_arguments", "error creating log interface to the log table, main store is nil")
	}
	return &SchemaVersion{storage: store}, nil
}

func (c *SchemaVersion) Search(ctx context.Context, opt *model.SearchOptions, filters any) ([]*model.SchemaVersion, errors.AppError) {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return nil, appErr
	}
	base := ApplyFiltersToBuilder(c.GetQueryBaseFromSearchOptions(opt), logFieldsMap, filters)
	query, _, _ := base.ToSql()
	wlog.Debug(query)
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

func (c *SchemaVersion) Insert(ctx context.Context, version *model.SchemaVersion) errors.AppError {
	db, appErr := c.storage.Database()
	if appErr != nil {
		return appErr
	}
	query := `insert into logger.schema_version(created_on, created_by, schema_id, object_data, version, note)
				VALUES ($1, $2, $3, $4, logger.get_next_version($3), $5);`
	wlog.Debug(query)
	_, err := db.ExecContext(ctx, query, version.CreatedOn, version.CreatedBy.Id, version.SchemaId, version.ObjectData, version.Note)
	if err != nil {
		return errors.NewInternalError("postgres.schema_version.insert.query.error", err.Error())
	}
	return nil
}

func (c *SchemaVersion) ScanRows(rows *sql.Rows) ([]*model.SchemaVersion, errors.AppError) {
	if rows == nil {
		return nil, errors.NewInternalError("postgres.schema_version.scan.check_args.rows_nil", "rows are nil")
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.NewInternalError("postgres.schema_version.scan.get_columns.error", err.Error())
	}
	var versions []*model.SchemaVersion

	for rows.Next() {
		var version model.SchemaVersion
		binds := make([]func(dst *model.SchemaVersion) interface{}, 0, 0)
		for _, v := range cols {
			switch v {
			case "id":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.Id })
			case "created_on":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.CreatedOn })
			case "created_by_id":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.CreatedBy.Id })
			case "created_by_name":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.CreatedBy.Name })
			case "schema_id":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.SchemaId })
			case "object_data":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.ObjectData })
			case "version":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.Version })
			case "note":
				binds = append(binds, func(dst *model.SchemaVersion) interface{} { return &dst.Note })
			default:
				panic("postgres.log.scan.get_columns.error: columns gotten from sql don't respond to model columns. Unknown column: " + v)
			}
		}
		bindFunc := func(binds []func(schemaVersion *model.SchemaVersion) any) []any {
			var fields []any
			for _, v := range binds {
				fields = append(fields, v(&version))
			}
			return fields
		}
		err = rows.Scan(bindFunc(binds)...)
		if err != nil {
			return nil, errors.NewInternalError("postgres.schema_version.scan.scan.error", err.Error())
		}
		versions = append(versions, &version)
	}
	if len(versions) == 0 {
		return nil, errors.NewBadRequestError("postgres.schema_version.scan.check_no_rows.error", sql.ErrNoRows.Error())
	}
	return versions, nil
}

func (c *SchemaVersion) GetQueryBaseFromSearchOptions(opt *model.SearchOptions) sq.SelectBuilder {
	var fields []string
	if opt == nil {
		return c.GetQueryBase(c.getFields())
	}
	for _, v := range opt.Fields {
		if columnName, ok := schemaVersionFieldsMap[v]; ok {
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
	if opt.Sort != "" {
		splitted := strings.Split(opt.Sort, ":")
		if len(splitted) == 2 {
			order := splitted[0]
			column := splitted[1]
			if column == "created_by" {
				column = "created_by_name"
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

func (c *SchemaVersion) GetQueryBase(fields []string) sq.SelectBuilder {
	base := sq.Select(fields...).
		From("logger.schema_version").
		JoinClause("LEFT JOIN directory.wbt_user ON wbt_user.id = schema_version.created_by").
		PlaceholderFormat(sq.Dollar)

	return base
}

func (c *SchemaVersion) getFields() []string {
	var fields []string
	for _, value := range schemaVersionFieldsMap {
		fields = append(fields, value)
	}
	return fields
}
