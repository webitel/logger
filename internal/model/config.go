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
	CreatedAt:       "created_at",
	CreatedBy:       "created_by",
	UpdatedAt:       "updated_at",
	UpdatedBy:       "updated_by",
	Enabled:         "enabled",
	DaysToStore:     "days_to_store",
	Period:          "period",
	NextUploadOn:    "next_upload_on",
	Storage:         "storage",
	DomainId:        "domain_id",
	Description:     "description",
	LastUploadedLog: "last_uploaded_log",
	LogsSize:        "logs_size",
	LogsCount:       "logs_count",
}
