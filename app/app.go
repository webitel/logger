package app

import (
	"database/sql"
	"strings"
	"webitel_logger/storage"

	errors "github.com/webitel/engine/model"
)

type App struct {
	storage storage.Storage
}

func New(store storage.Storage) (*App, errors.AppError) {
	if store == nil {
		return nil, errors.NewInternalError("app.app.new.check_arguments.fail", "store is nil")
	}
	return &App{storage: store}, nil
}

func IsErrNoRows(err errors.AppError) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}
