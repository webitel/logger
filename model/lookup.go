package model

type Lookup struct {
	Id   *NullInt    `json:"id,omitempty" db:"id"`
	Name *NullString `json:"name,omitempty" db:"name"`
}

func (l *Lookup) IsZero() bool {
	if l.Id.IsZero() && l.Name.IsZero() {
		return true
	}
	return false
}
