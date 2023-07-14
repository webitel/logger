package model

type Config struct {
	Id           int
	ObjectId     int
	Enabled      bool
	DaysToStore  int
	Period       string
	NextUploadOn int64
	StorageId    int
	DomainId     int
}
