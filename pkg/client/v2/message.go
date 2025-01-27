package client

import (
	"encoding/json"
	"time"
)

type RequiredFields struct {
	UserId int    `json:"userId,omitempty"`
	UserIp string `json:"userIp,omitempty"`
	Date   int64  `json:"date,omitempty"`
	Action string `json:"action,omitempty"`
}

type Record struct {
	Id       int64  `json:"id,omitempty"`
	NewState []byte `json:"newState,omitempty"`
}

type Message struct {
	Records        []*Record `json:"records,omitempty"`
	RequiredFields `json:"requiredFields"`
}

type MessageOpts func(*Message) error

func NewMessage(userId int64, userIp string, action Action, recordId int64, recordBody any) (*Message, error) {
	body, err := json.Marshal(recordBody)
	if err != nil {
		return nil, err
	}
	return &Message{
		RequiredFields: RequiredFields{
			UserId: int(userId),
			UserIp: userIp,
			Date:   time.Now().UnixMilli(),
			Action: action.String(),
		},
		Records: []*Record{{
			Id:       recordId,
			NewState: body,
		}},
	}, nil
}

func NewCreateMessage(userId int64, userIp string, recordId int64, recordBody any) (*Message, error) {
	body, err := json.Marshal(recordBody)
	if err != nil {
		return nil, err
	}
	return &Message{
		RequiredFields: RequiredFields{
			UserId: int(userId),
			UserIp: userIp,
			Date:   time.Now().UnixMilli(),
			Action: CreateAction.String(),
		},
		Records: []*Record{{
			Id:       recordId,
			NewState: body,
		}},
	}, nil
}

func NewUpdateMessage(userId int64, userIp string, recordId int64, recordBody any) (*Message, error) {
	body, err := json.Marshal(recordBody)
	if err != nil {
		return nil, err
	}
	return &Message{
		RequiredFields: RequiredFields{
			UserId: int(userId),
			UserIp: userIp,
			Date:   time.Now().UnixMilli(),
			Action: UpdateAction.String(),
		},
		Records: []*Record{{
			Id:       recordId,
			NewState: body,
		}},
	}, nil
}

func NewDeleteMessage(userId int64, userIp string, recordId int64) (*Message, error) {
	return &Message{
		RequiredFields: RequiredFields{
			UserId: int(userId),
			UserIp: userIp,
			Date:   time.Now().UnixMilli(),
			Action: DeleteAction.String(),
		},
		Records: []*Record{{
			Id: recordId,
		}},
	}, nil
}
