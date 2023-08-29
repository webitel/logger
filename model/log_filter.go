package model

//import (
//	"github.com/Masterminds/squirrel"
//	"time"
//)
//
//type LogFilter struct {
//	Action   string
//	DateFrom time.Time
//	DateTo   time.Time
//	UserIp   string
//	UserId   int
//	ConfigId int
//}
//
//func (f LogFilter) GetAction() (string, bool) {
//	if f.Action == "" {
//		return f.Action, false
//	}
//	return f.Action, true
//}
//
//func (f LogFilter) GetDateFrom() (time.Time, bool) {
//	if f.DateFrom.IsZero() {
//		return f.DateFrom, false
//	}
//	return f.DateFrom, true
//}
//func (f LogFilter) GetDateTo() (time.Time, bool) {
//	if f.DateTo.IsZero() {
//		return f.DateTo, false
//	}
//	return f.DateTo, true
//}
//func (f LogFilter) GetUserIp() (string, bool) {
//	if f.DateTo.IsZero() {
//		return f.UserIp, false
//	}
//	return f.UserIp, true
//}
//func (f LogFilter) GetUserId() (int, bool) {
//	if f.DateTo.IsZero() {
//		return f.UserId, false
//	}
//	return f.UserId, true
//}
//func (f LogFilter) GetConfigId() (int, bool) {
//	if f.DateTo.IsZero() {
//		return f.UserId, false
//	}
//	return f.UserId, true
//}
//
//type LogFilterBuilder struct {
//	action   string
//	dateFrom time.Time
//	dateTo   time.Time
//	userIp   string
//	userId   int
//	configId int
//
//	fields []string
//}
//
//func (l *LogFilterBuilder) Action(action string) *LogFilterBuilder {
//	l.action = action
//	return l
//}
//
//func (l *LogFilterBuilder) DateFrom(date time.Time) *LogFilterBuilder {
//	l.dateFrom = date
//	return l
//}
//func (l *LogFilterBuilder) DateTo(date time.Time) *LogFilterBuilder {
//	l.dateTo = date
//	return l
//}
//func (l *LogFilterBuilder) UserIp(userIp string) *LogFilterBuilder {
//	l.userIp = userIp
//	return l
//}
//func (l *LogFilterBuilder) UserId(userId int) *LogFilterBuilder {
//	l.userId = userId
//	return l
//}
//
//func NewLogFilterBuilder() *LogFilterBuilder {
//	return &LogFilterBuilder{}
//}
//
//func (l *LogFilterBuilder) Build() *LogFilter {
//	return &LogFilter{
//		Action:   l.action,
//		DateFrom: l.dateFrom,
//		DateTo:   l.dateTo,
//		UserIp:   l.userIp,
//		UserId:   l.userId,
//		ConfigId: l.configId,
//	}
//}

type Filter struct {
	Column         string
	Value          any
	ComparisonType Comparison
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
