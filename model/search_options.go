package model

type SearchOptions struct {
	Page   int    `json:"page,omitempty"`
	Size   int    `json:"size,omitempty"`
	Search string `json:"q,omitempty"`

	Sort   string   `json:"sort,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

func (s *SearchOptions) GetSize() int32 {
	return int32(s.Size)
}
