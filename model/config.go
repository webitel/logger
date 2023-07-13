package model

import "time"

type Config struct {
	Id           int
	ObjectId     int
	Enabled      bool
	DaysToStore  int
	Period       string
	NextUploadOn time.Time
	StorageId    int
	DomainId     int
}
