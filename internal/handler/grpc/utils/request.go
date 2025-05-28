package utils

import (
	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/model"
)

func UnmarshalLookup[K model.Lookup](lp *proto.Lookup, lookup K) K {
	if lp == nil {
		var res K
		return res
	}
	lookup.SetId(int(lp.Id))
	lookup.SetName(lp.Name)
	return lookup
}
