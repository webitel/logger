package auth

import (
	"context"
	"github.com/golang/groupcache/singleflight"
	authproto "github.com/webitel/logger/api/webitel-go"
	"github.com/webitel/logger/internal/auth/model"
	errors "github.com/webitel/logger/internal/model"
	"google.golang.org/grpc"
)

type AuthManager interface {
	Authorize(ctx context.Context, token string) (*model.Session, errors.AppError)
	AuthorizeFromContext(ctx context.Context) (*model.Session, errors.AppError)
}

type AuthorizationClient struct {
	Client     authproto.AuthClient
	Group      singleflight.Group
	Connection *grpc.ClientConn
}

func NewAuthorizationClient(conn *grpc.ClientConn) (*AuthorizationClient, errors.AppError) {
	if conn == nil {
		return nil, errors.NewInternalError("auth.manager.new_auth_client.validate_params.connection", "invalid GRPC connection")
	}
	return &AuthorizationClient{
		Client:     authproto.NewAuthClient(conn),
		Group:      singleflight.Group{},
		Connection: conn,
	}, nil
}

func (c *AuthorizationClient) UserInfo(ctx context.Context, token string) (*model.Session, errors.AppError) {
	interfacedSession, err := c.Group.Do(token, func() (interface{}, error) {
		info, err := c.Client.UserInfo(ctx, &authproto.UserinfoRequest{AccessToken: token})
		if err != nil {
			return nil, err
		}
		return ConstructSessionFromUserInfo(info), nil
	})
	if err != nil {
		return nil, errors.NewUnauthorizedError("auth.manager.user_info.do_request.error", err.Error())
	}
	return interfacedSession.(*model.Session), nil
}

func ConstructSessionFromUserInfo(userinfo *authproto.Userinfo) *model.Session {
	session := &model.Session{
		User: &model.AuthorizedUser{
			Id:        userinfo.UserId,
			Name:      userinfo.Name,
			Username:  userinfo.Username,
			Extension: userinfo.Extension,
		},
		ExpiresAt: userinfo.ExpiresAt,
		DomainId:  userinfo.Dc,
	}
	for i, permission := range userinfo.Permissions {
		if i == 0 {
			session.Permissions = make([]*model.Permission, 0)
		}
		session.Permissions = append(session.Permissions, &model.Permission{
			Id:   permission.GetId(),
			Name: permission.GetName(),
		})
	}
	for i, scope := range userinfo.Scope {
		if i == 0 {
			session.Scope = make([]*model.Scope, 0)
		}
		session.Scope = append(session.Scope, &model.Scope{
			Id:     scope.GetId(),
			Name:   scope.GetName(),
			Abac:   scope.Abac,
			Obac:   scope.Obac,
			Rbac:   scope.Rbac,
			Class:  scope.Class,
			Access: scope.Access,
		})
	}

	for i, role := range userinfo.Roles {
		if i == 0 {
			session.Roles = make([]*model.Role, 0)
		}
		session.Roles = append(session.Roles, &model.Role{
			Id:   role.GetId(),
			Name: role.GetName(),
		})
	}
	return session
}
