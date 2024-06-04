package model

type BrokerLogMessage struct {
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
	Id       int64     `json:"id,omitempty"`
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
