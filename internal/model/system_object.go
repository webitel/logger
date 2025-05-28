package model

type SystemObject struct {
	Id   *int    `db:"id"`
	Name *string `db:"name"`
}

func (a *SystemObject) SetId(id int) {
	if a == nil {
		return
	}
	a.Id = &id
}

func (a *SystemObject) SetName(name string) {
	if a == nil {
		return
	}
	a.Name = &name
}

func (a *SystemObject) GetId() *int {
	if a == nil {
		return nil
	}
	return a.Id
}

func (a *SystemObject) GetName() *string {
	if a == nil {
		return nil
	}
	return a.Name
}
