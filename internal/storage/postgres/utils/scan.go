package utils

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/logger/internal/model"
)

func ScanRows[K any](rows *pgx.Rows, getPlanByColumns func([]string) []func(*K) any) ([]*K, error) {
	if rows == nil {
		return nil, errors.New("rows is nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var res []*K
	for rows.Next() {
		var (
			entity       K
			destinations []any
		)
		for _, bind := range getPlanByColumns(columns) {
			destinations = append(destinations, bind(&entity))

		}
		err = rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}
		res = append(res, &entity)
	}
	return res, nil
}

func ScanRow[K any](row *pgx.Rows, getPlanByColumns func([]string) []func(*K) any) (*K, error) {
	res, err := ScanRows(row, getPlanByColumns)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, sql.ErrNoRows
	}
	return res[0], nil
}

func ScanLookupId(lookup **model.Lookup) any {
	orig := *lookup
	if orig == nil {
		orig = new(model.Lookup)
	}
	return &orig.Id
}
func ScanLookupName(lookup **model.Lookup) any {
	orig := *lookup
	if orig == nil {
		orig = new(model.Lookup)
	}
	return &orig.Name
}
