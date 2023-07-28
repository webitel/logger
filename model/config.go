package model

type Config struct {
	Id           int
	ObjectId     int
	CreatedAt    NullTime
	CreatedBy    int
	UpdatedAt    NullTime
	UpdatedBy    int
	Enabled      bool
	DaysToStore  int
	Period       string
	NextUploadOn NullTime
	StorageId    int
	DomainId     int
}
