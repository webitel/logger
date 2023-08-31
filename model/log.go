package model

type Log struct {
	Id       int
	Action   string
	Date     NullTime
	User     Lookup
	Object   Lookup
	UserIp   string
	RecordId int
	NewState []byte
	ConfigId int
}
