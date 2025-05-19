package utils

import "github.com/webitel/logger/internal/model"

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
func ConvertToOutputBulk[C any, K any](items []C, convertFunc func(C) (K, model.AppError)) ([]K, model.AppError) {
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

// C type of input, K type of output
func CalculateListResultMetadata[C any, K any](s Lister, items []C, convertFunc func(C) (K, model.AppError)) (bool, []K, model.AppError) {
	var (
		result []K
		err    model.AppError
	)
	next, filteredInput := GetListResult[C](s, items)
	result, err = ConvertToOutputBulk[C, K](filteredInput, convertFunc)
	if err != nil {
		return false, nil, err
	}
	return next, result, nil
}
