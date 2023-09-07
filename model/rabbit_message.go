package model

type RabbitMessage struct {
	//ObjectId int    `json:"objectId,omitempty"`
	NewState []byte `json:"newState,omitempty"`
	UserId   int    `json:"userId,omitempty"`
	UserIp   string `json:"userIp,omitempty"`
	Action   string `json:"action,omitempty"`
	Date     int64  `json:"date,omitempty"`
	//DomainId int    `json:"domainId,omitempty"`
	RecordId int64 `json:"recordId,omitempty"`
}
