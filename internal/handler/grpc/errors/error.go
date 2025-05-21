package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrDatabase = NewInternalError("api.process.database.error", "database error")
	ErrInternal = NewInternalError("api.process.request.error", "an unpredictable error occurred while processing request")
)

type AppError interface {
	SetStatusCode(int) AppError
	GetStatusCode() int
	SetDetailedError(string)
	GetDetailedError() string
	SetRequestId(string)
	GetRequestId() string
	GetId() string

	Error() string
	ToJson() string
	String() string
}

type ApplicationError struct {
	Id            string `json:"id"`
	Where         string `json:"where,omitempty"`
	Status        string `json:"status"`
	DetailedError string `json:"detail"`
	RequestId     string `json:"request_id,omitempty"`
	StatusCode    int    `json:"code,omitempty"`
}

func (err *ApplicationError) SetStatusCode(code int) AppError {
	err.StatusCode = code
	err.Status = http.StatusText(err.StatusCode)
	return err
}

func (err *ApplicationError) GetStatusCode() int {
	return err.StatusCode
}

func (err *ApplicationError) Error() string {
	var where string
	if err.Where != "" {
		where = err.Where + ": "
	}
	return fmt.Sprintf("%s%s, %s", where, err.Status, err.DetailedError)
}

func (err *ApplicationError) SetDetailedError(details string) {
	err.DetailedError = details
}

func (err *ApplicationError) GetDetailedError() string {
	return err.DetailedError
}

func (err *ApplicationError) SetRequestId(id string) {
	err.RequestId = id
}

func (err *ApplicationError) GetRequestId() string {
	return err.RequestId
}

func (err *ApplicationError) GetId() string {
	return err.Id
}

func (err *ApplicationError) ToJson() string {
	b, _ := json.Marshal(err)
	return string(b)
}

func (err *ApplicationError) String() string {
	if err.Id == err.Status && err.DetailedError != "" {
		return err.DetailedError
	}
	return err.Status
}

// Error constructors
func NewInternalError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusInternalServerError)
}

func NewNotFoundError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusNotFound)
}

func NewBadRequestError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusBadRequest)
}

func NewForbiddenError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusForbidden)
}

func newAppError(id string, details string) AppError {
	return &ApplicationError{Id: id, Status: id, DetailedError: details}
}
