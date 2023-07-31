package model

type Lookup struct {
	Id   *NullInt    `json:"id,omitempty"`
	Name *NullString `json:"name,omitempty"`
}

func (l *Lookup) IsZero() bool {
	if l.Id.IsZero() && l.Name.IsZero() {
		return true
	}
	return false
}
