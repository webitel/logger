package model

type Filter struct {
	Column         string
	Value          any
	ComparisonType Comparison
}

type FilterBunch struct {
	Bunch []*Filter
	ConnectionType
}

type FilterArray struct {
	Filters    []*FilterBunch
	Connection ConnectionType
}

type Comparison int64

const (
	Equal              Comparison = 0
	GreaterThan        Comparison = 1
	GreaterThanOrEqual Comparison = 2
	LessThan           Comparison = 3
	LessThanOrEqual    Comparison = 4
	NotEqual           Comparison = 5
	Like               Comparison = 6
	ILike              Comparison = 7
)

type ConnectionType int64

const (
	AND ConnectionType = 0
	OR  ConnectionType = 1
)
