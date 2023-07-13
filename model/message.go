package model

type Message struct {
	ObjectId int    `json:"object,omitempty"`
	NewState string `json:"newState,omitempty"`
	UserId   int    `json:"userId,omitempty"`
	UserIp   string `json:"userIp,omitempty"`
	Action   string `json:"action,omitempty"`
	Date     string `json:"date,omitempty"`
	DomainId int    `json:"domainId,omitempty"`
}
