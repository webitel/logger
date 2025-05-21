package model

import (
	guid "github.com/google/uuid"
	"time"
)

type BrokerRecordLogMessage struct {
	Records                        []*LogEntity `json:"records,omitempty"`
	BrokerLogMessageRequiredFields `json:"requiredFields"`
}

type BrokerLogMessageRequiredFields struct {
	UserId int    `json:"userId,omitempty"`
	UserIp string `json:"userIp,omitempty"`
	Action string `json:"action,omitempty"`
	Date   int64  `json:"date,omitempty"`
}

type LogEntity struct {
	Id       int       `json:"id,omitempty"`
	NewState BytesJSON `json:"newState,omitempty"`
}

type BytesJSON struct {
	Body []byte
}

func (b *BytesJSON) GetBody() []byte {
	return b.Body
}

func (b *BytesJSON) UnmarshalJSON(input []byte) error {
	b.Body = input
	return nil
}

type BrokerLoginMessage struct {
	AuthType string  `json:"type,omitempty"`   // Authentication type (most likely "password")
	Agent    string  `json:"agent,omitempty"`  // User-Agent
	Date     int64   `json:"date,omitempty"`   // Date of event
	IsNew    bool    `json:"isNew,omitempty"`  // True - if user logged-in
	From     string  `json:"from,omitempty"`   // IP address of user
	Login    *Login  `json:"login,omitempty"`  // Context of login
	Status   *Status `json:"status,omitempty"` // Login operation status
}

func (m *BrokerLoginMessage) ConvertToDatabaseModel() (*LoginAttempt, error) {
	var (
		success       bool
		databaseModel LoginAttempt
		authType      string
	)
	if m.Status != nil {
		success = false
		databaseModel.Details = NewNullString(m.Status.Detail)
	} else {
		success = true
	}
	if user := m.Login.User; user != nil {
		if user.Id != 0 {
			id, err := NewNullInt(user.Id)
			if err != nil {
				return nil, err
			}
			databaseModel.UserId = id
		}
		databaseModel.UserName = user.Username

	}
	if domain := m.Login.Domain; domain != nil {
		if domain.Id != 0 {
			id, err := NewNullInt(domain.Id)
			if err != nil {
				return nil, err
			}
			databaseModel.DomainId = id
		}
		databaseModel.DomainName = domain.Name

	}
	authType = m.AuthType
	if authType == "" {
		authType = "password"
	}
	databaseModel.AuthType = authType
	databaseModel.UserAgent = m.Agent
	databaseModel.Date = time.UnixMilli(m.Date)
	databaseModel.UserIp = m.From
	databaseModel.Success = success

	return &databaseModel, nil
}

type Login struct {
	Id        *guid.UUID        `json:"id,omitempty"`
	CreatedAt int64             `json:"created_at,omitempty"`
	ExpiresAt int64             `json:"expires_at,omitempty"`
	MaxAge    int64             `json:"max_age,omitempty"`
	Context   map[string]string `json:"context,omitempty"`
	Domain    *Domain           `json:"domain,omitempty"`
	User      *User             `json:"user,omitempty"`
}

type User struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	Extension string `json:"extension,omitempty"`
}

type Status struct {
	Id     string `json:"id,omitempty"`
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type Domain struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
