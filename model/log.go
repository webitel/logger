package model

import (
	"github.com/webitel/engine/model"
	"time"
)

type Log struct {
	Id       int
	Action   string
	Date     time.Time
	User     model.Lookup
	UserIp   string
	RecordId int
	NewState string
	ConfigId int
}
