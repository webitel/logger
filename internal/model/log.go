package model

import (
	"time"
)

type Log struct {
	*Object
	*Author
	*Record
	Id       int        `json:"id,omitempty"`
	Action   string     `json:"action,omitempty"`
	Date     *time.Time `json:"date,omitempty"`
	UserIp   *string    `json:"user_ip,omitempty"`
	NewState []byte     `json:"new_state,omitempty"`
	ConfigId int        `json:"config_id,omitempty"`
}

type Record struct {
	Id   *int    `db:"record_id"`
	Name *string `db:"record_name"`
}

func (r *Record) GetId() *int {
	if r == nil {
		return nil
	}
	return r.Id
}

func (r *Record) GetName() *string {
	if r == nil {
		return nil
	}
	return r.Name
}

func (r *Record) SetId(id int) {
	if r == nil {
		return
	}
	r.Id = &id
}

func (r *Record) SetName(name string) {
	if r == nil {
		return
	}
	r.Name = &name
}

// Front-end fields
var LogFields = struct {
	Id       string
	Action   string
	Date     string
	User     string
	Object   string
	UserIp   string
	Record   string
	NewState string
	ConfigId string
}{
	Id:       "id",
	Action:   "action",
	Date:     "date",
	User:     "user",
	Object:   "object",
	UserIp:   "user_ip",
	Record:   "record",
	NewState: "new_state",
	ConfigId: "config_id",
}

type LogFilters struct {
	Id       []int64
	Action   []string
	DateFrom *time.Time
	DateTo   *time.Time
	User     []int64
	Object   []int64
	UserIp   []string
	Record   []int64
	ConfigId []int64
}

func (l *LogFilters) ExtractFilters() *FilterNode {
	main := &FilterNode{
		Nodes:      make([]any, 0),
		Connection: AND,
	}
	res := GetOrFiltersFromArray[int64](l.Id, LogFields.Id, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	res = GetOrFiltersFromArray[string](l.Action, LogFields.Action, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	res = GetOrFiltersFromArray[string](l.UserIp, LogFields.UserIp, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	res = GetOrFiltersFromArray[int64](l.User, LogFields.User, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	res = GetOrFiltersFromArray[int64](l.Record, LogFields.Record, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	res = GetOrFiltersFromArray[int64](l.ConfigId, LogFields.ConfigId, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	res = GetOrFiltersFromArray[int64](l.Object, LogFields.Object, Equal)
	if res != nil {
		main.Nodes = append(main.Nodes, res)
	}
	if l.DateFrom != nil {
		main.Nodes = append(main.Nodes, &Filter{
			Column:         LogFields.Date,
			Value:          l.DateFrom,
			ComparisonType: GreaterThanOrEqual,
		})
	}
	if l.DateTo != nil {
		main.Nodes = append(main.Nodes, &Filter{
			Column:         LogFields.Date,
			Value:          l.DateTo,
			ComparisonType: LessThanOrEqual,
		})
	}
	return main
}

func GetOrFiltersFromArray[C any](in []C, fieldName string, comparison Comparison) any {
	if len(in) > 0 {
		if len(in) == 1 {
			return &Filter{
				Column:         fieldName,
				Value:          in[0],
				ComparisonType: comparison,
			}
		}
		sub := &FilterNode{Nodes: make([]any, 0), Connection: OR}
		for _, action := range in {
			sub.Nodes = append(sub.Nodes, &Filter{
				Column:         fieldName,
				Value:          action,
				ComparisonType: comparison,
			})
		}
		return sub
	}
	return nil
}
