package model

type Lookup interface {
	GetId() *int
	GetName() *string
	SetId(id int)
	SetName(name string)
}
