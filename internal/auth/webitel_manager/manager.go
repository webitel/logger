package webitel_manager

import (
	"context"
	iface "github.com/webitel/logger/internal/auth"
	authmodel "github.com/webitel/logger/internal/auth/model"
	errors "github.com/webitel/logger/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	RequestContextName = "grpc_ctx"
	AuthTokenName      = "X-Webitel-Access"
)

var _ iface.AuthManager = &WebitelAppAuthManager{}

type WebitelAppAuthManager struct {
	client *iface.AuthorizationClient
}

func NewWebitelAppAuthManager(conn *grpc.ClientConn) (iface.AuthManager, errors.AppError) {
	cli, err := iface.NewAuthorizationClient(conn)
	if err != nil {
		return nil, err
	}
	manager := &WebitelAppAuthManager{client: cli}

	return manager, nil
}

func (i *WebitelAppAuthManager) AuthorizeFromContext(ctx context.Context) (*authmodel.Session, errors.AppError) {
	var token []string
	var info metadata.MD
	var ok bool

	v := ctx.Value(RequestContextName)
	info, ok = v.(metadata.MD)

	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return nil, errors.NewForbiddenError("app.grpc.get_context", "Not found")
	} else {
		token = info.Get(AuthTokenName)
	}
	newContext := metadata.NewOutgoingContext(ctx, info)
	if len(token) < 1 {
		return nil, errors.NewInternalError("webitel_manager.authorize_from_from_context.search_token.not_found", "token not found")
	}
	return i.Authorize(newContext, token[0])

}

func (i *WebitelAppAuthManager) Authorize(ctx context.Context, token string) (*authmodel.Session, errors.AppError) {
	return i.client.UserInfo(ctx, token)
}
