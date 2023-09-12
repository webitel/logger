package model

type Log struct {
	Id       int
	Action   string
	Date     NullTime
	User     Lookup
	Object   Lookup
	UserIp   string
	Record   Lookup
	NewState []byte
	ConfigId int
}
