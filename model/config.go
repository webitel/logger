package model

type Config struct {
	Id           int
	Object       Lookup
	CreatedAt    NullTime
	CreatedBy    int
	UpdatedAt    NullTime
	UpdatedBy    NullInt
	Enabled      bool
	DaysToStore  int
	Period       int
	NextUploadOn NullTime
	Storage      Lookup
	DomainId     int
}
