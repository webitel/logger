package model

type Log struct {
	Id       int
	Action   string
	Date     int64
	UserId   int
	UserIp   string
	ObjectId int
	NewState string
	DomainId int
}
