package utils

import (
	"database/sql"
	"errors"
)

func ScanRows[K any](rows *sql.Rows, getPlanByColumns func([]string) []func(*K) any) ([]*K, error) {
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

func ScanRow[K any](row *sql.Rows, getPlanByColumns func([]string) []func(*K) any) (*K, error) {
	res, err := ScanRows(row, getPlanByColumns)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, sql.ErrNoRows
	}
	return res[0], nil
}
