package model

import (
	"fmt"
)

const (
	DefaultPageSize = 40
	DefaultPage     = 1
	MaxPageSize     = 40000
)

type Searcher interface {
	GetPage() int32
	GetSize() int32
	GetQ() string
	GetSort() string
	GetFields() []string
}

func ConvertSort(in string) string {
	if len(in) < 2 || (in[0] != '+' && in[0] != '-') {
		return ""
	}
	if in[0] == '+' {
		return fmt.Sprintf("%s:%s", "ASC", in[1:])
	} else {
		return fmt.Sprintf("%s:%s", "DESC", in[1:])
	}
}

type SearchOptions struct {
	Page   int    `json:"page,omitempty"`
	Size   int    `json:"size,omitempty"`
	Search string `json:"q,omitempty"`

	Sort   string   `json:"sort,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

func (s *SearchOptions) GetSize() int {
	return s.Size
}
