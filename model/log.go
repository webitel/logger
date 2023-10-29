package model

type Log struct {
	Id       int      `json:"id,omitempty"`
	Action   string   `json:"action,omitempty"`
	Date     NullTime `json:"date,omitempty"`
	User     Lookup   `json:"user,omitempty"`
	Object   Lookup   `json:"object,omitempty"`
	UserIp   string   `json:"user_ip,omitempty"`
	Record   Lookup   `json:"record,omitempty"`
	NewState []byte   `json:"new_state,omitempty"`
	ConfigId int      `json:"config_id,omitempty"`
}
