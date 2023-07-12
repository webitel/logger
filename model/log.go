package model

type Log struct {
	Action   string
	Date     string
	UserId   int
	UserIp   string
	ObjectId int
	NewState string
	DomainId int
}
