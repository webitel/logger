package model

type Config struct {
	Id              int
	Object          Lookup
	CreatedAt       NullTime
	CreatedBy       int
	UpdatedAt       NullTime
	UpdatedBy       NullInt
	Enabled         bool
	DaysToStore     int
	Period          int
	NextUploadOn    NullTime
	Storage         Lookup
	DomainId        int64
	Description     NullString
	LastUploadedLog NullInt
	LogsSize        NullString
	LogsCount       NullInt
}

// Front-end fields
var ConfigFields = struct {
	Id              string
	Object          string
	CreatedAt       string
	CreatedBy       string
	UpdatedAt       string
	UpdatedBy       string
	Enabled         string
	DaysToStore     string
	Period          string
	NextUploadOn    string
	Storage         string
	DomainId        string
	Description     string
	LastUploadedLog string
	LogsSize        string
	LogsCount       string
}{
	Id:              "id",
	Object:          "object",
	CreatedAt:       "createdAt",
	CreatedBy:       "createdBy",
	UpdatedAt:       "updatedAt",
	UpdatedBy:       "updatedBy",
	Enabled:         "enabled",
	DaysToStore:     "daysToStore",
	Period:          "period",
	NextUploadOn:    "nextUploadOn",
	Storage:         "storage",
	DomainId:        "domainId",
	Description:     "description",
	LastUploadedLog: "lastUploadedLog",
	LogsSize:        "logsSize",
	LogsCount:       "logsCount",
}
