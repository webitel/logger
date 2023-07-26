package cache

import (
	"context"

	errors "github.com/webitel/engine/model"
)

type CacheValue struct {
	value any
}

type CacheStore interface {
	Get(ctx context.Context, key string) (*CacheValue, errors.AppError)
	Set(ctx context.Context, key string, value any, expiresAfter int64) errors.AppError
	Delete(ctx context.Context, key string) errors.AppError
}

func (v *CacheValue) String() (string, errors.AppError) {
	if v.value == nil {
		return "", errors.NewInternalError("cache.cache.string.check_value.value_nil", "value is nil")
	}
	if v, ok := v.value.(string); ok {
		return v, nil
	} else {
		return "", errors.NewInternalError("cache.cache.string.convert_value.error", "unable to convert value")
	}

}

func (v *CacheValue) Raw() any {
	return v.value
}
func (v *CacheValue) Set(value any) errors.AppError {
	if value != nil {
		v.value = value
		return nil
	} else {
		return errors.NewInternalError("cache.cache.string.set.check_arguments.error", "accepted value is nil")
	}
}

func NewCacheValue(value any) (*CacheValue, errors.AppError) {
	var cv CacheValue
	err := cv.Set(value)
	if err != nil {
		return nil, err
	}
	return &cv, err
}
