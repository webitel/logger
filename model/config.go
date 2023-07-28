package model

type Config struct {
	Id           int
	ObjectId     int
	CreatedAt    NullTime
	CreatedBy    int
	UpdatedAt    NullTime
	UpdatedBy    NullInt
	Enabled      bool
	DaysToStore  int
	Period       string
	NextUploadOn NullTime
	StorageId    int
	DomainId     int
}
