package model

import "time"

type SchemaVersion struct {
	Id         int64       `json:"id,omitempty"`
	SchemaId   int64       `json:"schemaId,omitempty"`
	CreatedOn  time.Time   `json:"createdOn"`
	CreatedBy  Lookup      `json:"createdBy"`
	ObjectData []byte      `json:"objectData,omitempty"`
	Version    int64       `json:"version,omitempty"`
	Note       *NullString `json:"note,omitempty"`
}
