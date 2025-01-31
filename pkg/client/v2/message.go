package client

import (
	"time"
)

type RequiredFields struct {
	UserId int    `json:"userId,omitempty"`
	UserIp string `json:"userIp,omitempty"`
	Date   int64  `json:"date,omitempty"`
	Action string `json:"action,omitempty"`
}

type Record struct {
	Id       int64 `json:"id,omitempty"`
	NewState any   `json:"newState,omitempty"`
}

type Message struct {
	Records        []*Record `json:"records,omitempty"`
	RequiredFields `json:"requiredFields"`
}

type MessageOpts func(*Message) error

func NewMessage(userId int64, userIp string, action Action, recordId int64, recordBody any) (*Message, error) {
	return &Message{
		RequiredFields: RequiredFields{
			UserId: int(userId),
			UserIp: userIp,
			Date:   time.Now().Unix(),
			Action: action.String(),
		},
		Records: []*Record{{
			Id:       recordId,
			NewState: recordBody,
		}},
	}, nil
}

func NewCreateMessage(userId int64, userIp string, recordId int64, recordBody any) (*Message, error) {
	return NewMessage(userId, userIp, CreateAction, recordId, recordBody)
}

func NewUpdateMessage(userId int64, userIp string, recordId int64, recordBody any) (*Message, error) {
	return NewMessage(userId, userIp, UpdateAction, recordId, recordBody)
}

func NewDeleteMessage(userId int64, userIp string, recordId int64) (*Message, error) {
	return NewMessage(userId, userIp, DeleteAction, recordId, nil)
}
