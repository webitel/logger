package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthError interface {
	SetStatusCode(int) AuthError
	GetStatusCode() int
	SetDetailedError(string)
	GetDetailedError() string
	SetRequestId(string)
	GetRequestId() string
	GetId() string

	Error() string
	ToJson() string
	String() string
	AuthError()
}

type AuthorizationError struct {
	params        map[string]interface{}
	Id            string `json:"id"`
	Where         string `json:"where,omitempty"`
	Status        string `json:"status"`
	DetailedError string `json:"detail"`
	RequestId     string `json:"request_id,omitempty"`
	StatusCode    int    `json:"code,omitempty"`
}

func (err *AuthorizationError) AuthError() {}

func (err *AuthorizationError) SetStatusCode(code int) *AuthorizationError {
	err.StatusCode = code
	err.Status = http.StatusText(err.StatusCode)
	return err
}

func (err *AuthorizationError) GetStatusCode() int {
	return err.StatusCode
}

func (err *AuthorizationError) Error() string {
	return fmt.Sprintf("AuthError [%s]: %s, %s", err.Id, err.Status, err.DetailedError)
}

func (err *AuthorizationError) SetDetailedError(details string) {
	err.DetailedError = details
}

func (err *AuthorizationError) GetDetailedError() string {
	return err.DetailedError
}

func (err *AuthorizationError) SetRequestId(id string) {
	err.RequestId = id
}

func (err *AuthorizationError) GetRequestId() string {
	return err.RequestId
}

func (err *AuthorizationError) GetId() string {
	return err.Id
}

func (err *AuthorizationError) ToJson() string {
	b, _ := json.Marshal(err)
	return string(b)
}

func (err *AuthorizationError) String() string {
	if err.Id == err.Status && err.DetailedError != "" {
		return err.DetailedError
	}
	return err.Status
}

func NewUnauthorizedError(id, details string) *AuthorizationError {
	return newAuthError(id, details).SetStatusCode(http.StatusUnauthorized)
}

func NewForbiddenError(id, details string) *AuthorizationError {
	return newAuthError(id, details).SetStatusCode(http.StatusForbidden)
}

func newAuthError(id string, details string) *AuthorizationError {
	return &AuthorizationError{Id: id, Status: id, DetailedError: details}
}
