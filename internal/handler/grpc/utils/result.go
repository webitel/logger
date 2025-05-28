package utils

import (
	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/handler/grpc/errors"
	"github.com/webitel/logger/internal/model"
)

type Lister interface {
	GetSize() int32
}

// C type of items to filter
func GetListResult[C any](s Lister, items []C) (bool, []C) {
	if int32(len(items)-1) == s.GetSize() {
		return true, items[0 : len(items)-1]
	}
	return false, items
}

// C type of input, K type of output
func ConvertToOutputBulk[C any, K any](items []C, convertFunc func(C) (K, errors.AppError)) ([]K, errors.AppError) {
	var result []K
	for _, item := range items {
		out, err := convertFunc(item)
		if err != nil {
			return nil, err
		}
		result = append(result, out)
	}
	return result, nil
}

func ResolvePaging[C any](size int, items []C) (result []C, next bool) {
	if len(items) >= size {
		return items[0:size], true
	}
	return items, false
}

func MarshalLookup(lookup model.Lookup) *proto.Lookup {
	if lookup == nil {
		return nil
	}
	m := &proto.Lookup{}
	if id := lookup.GetId(); id != nil {
		m.Id = int32(*id)
	}
	if name := lookup.GetName(); name != nil {
		m.Name = *name
	}
	return m
}
