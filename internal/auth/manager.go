package auth

import (
	"context"
)

type Manager interface {
	AuthorizeFromContext(ctx context.Context, mainObjClass string, mainAccessMode AccessMode) (Auther, error)
}
