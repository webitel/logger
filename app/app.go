package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"webitel_logger/model"
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

func ExtractSearchOptions(t any) (*model.SearchOptions, errors.AppError) {
	var res model.SearchOptions
	b, err := json.Marshal(t)
	if err != nil {
		return nil, errors.NewBadRequestError("app.app.extract_search_options.marshal.error", err.Error())
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.NewInternalError("app.app.extract_search_options.unmarshal.error", err.Error())
	}
	if res.Sort != "" {
		res.Sort = ConvertSort(res.Sort)
	}
	return &res, nil
}

func ConvertSort(in string) string {
	if len(in) < 2 {
		return ""
	}
	if in[0] == '+' {
		return fmt.Sprintf("%s:%s", "ASC", in[1:])
	} else {
		return fmt.Sprintf("%s:%s", "DESC", in[1:])
	}
}
