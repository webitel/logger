package model

import "time"

type Config struct {
	*Storage
	*Author
	*Editor
	*Object
	Id                int
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	Enabled           bool
	DaysToStore       int
	Period            int
	NextUploadOn      *time.Time
	StorageId         *int
	DomainId          int
	Description       *string
	LastUploadedLogId *int
	LogsSize          *string
	LogsCount         *int
}

type Storage struct {
	Id   *int    `db:"storage_id"`
	Name *string `db:"storage_name"`
}

func (a *Storage) SetId(id int) {
	if a == nil {
		return
	}
	a.Id = &id
}

func (a *Storage) SetName(name string) {
	if a == nil {
		return
	}
	a.Name = &name
}

func (a *Storage) GetId() *int {
	if a == nil {
		return nil
	}
	return a.Id
}

func (a *Storage) GetName() *string {
	if a == nil {
		return nil
	}
	return a.Name
}

type Object struct {
	Id   *int    `db:"object_id"`
	Name *string `db:"object_name"`
}

func (a *Object) SetId(id int) {
	if a == nil {
		return
	}
	a.Id = &id
}

func (a *Object) SetName(name string) {
	if a == nil {
		return
	}
	a.Name = &name
}

func (a *Object) GetId() *int {
	if a == nil {
		return nil
	}
	return a.Id
}

func (a *Object) GetName() *string {
	if a == nil {
		return nil
	}
	return a.Name
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

func GetConfigFields() []string {
	return []string{
		ConfigFields.Id,
		ConfigFields.Object,
		ConfigFields.CreatedAt,
		ConfigFields.CreatedBy,
		ConfigFields.UpdatedAt,
		ConfigFields.UpdatedBy,
		ConfigFields.Enabled,
		ConfigFields.DaysToStore,
		ConfigFields.Period,
		ConfigFields.NextUploadOn,
		ConfigFields.Storage,
		ConfigFields.DomainId,
		ConfigFields.Description,
		ConfigFields.LastUploadedLog,
		ConfigFields.LogsSize,
		ConfigFields.LogsCount,
	}
}
