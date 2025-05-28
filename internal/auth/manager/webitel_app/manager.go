package webitel_app

import (
	"context"
	proto "github.com/webitel/logger/api/webitel-go"
	"github.com/webitel/logger/internal/auth"
	autherror "github.com/webitel/logger/internal/auth/errors"
	session "github.com/webitel/logger/internal/auth/session/user_session"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

var _ auth.Manager = &Manager{}

type Manager struct {
	Client     proto.AuthClient
	Group      singleflight.Group
	Connection *grpc.ClientConn
}

func New(conn *grpc.ClientConn) (*Manager, error) {
	return &Manager{Client: proto.NewAuthClient(conn), Group: singleflight.Group{}, Connection: conn}, nil
}

func (i *Manager) AuthorizeFromContext(ctx context.Context, mainObjClassName string, mainAccessMode auth.AccessMode) (auth.Auther, error) {
	var token []string
	var info metadata.MD
	var ok bool

	v := ctx.Value(session.RequestContextName)
	info, ok = v.(metadata.MD)

	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return nil, autherror.NewUnauthorizedError("internal.grpc.get_context", "Not found")
	} else {
		token = info.Get(session.AuthTokenName)
	}
	newContext := metadata.NewOutgoingContext(ctx, info)
	if len(token) < 1 {
		return nil, autherror.NewUnauthorizedError("webitel_manager.authorize_from_from_context.search_token.not_found", "token not found")
	}
	userToken := token[0]
	sess, err, _ := i.Group.Do(userToken, func() (interface{}, error) {
		return i.Client.UserInfo(newContext, nil)
	})
	if err != nil {
		return nil, autherror.NewUnauthorizedError("webitel_manager.authorize_from_from_context.user_info.err", err.Error())
	}
	auther := ConstructSessionFromUserInfo(sess.(*proto.Userinfo), mainObjClassName, mainAccessMode, getClientIp(ctx))
	if auther.IsExpired() {
		return nil, autherror.NewUnauthorizedError("webitel_manager.authorize_from_from_context.user_info.err", "session expired")
	}
	return auther, nil

}

func ConstructSessionFromUserInfo(userinfo *proto.Userinfo, mainObjClass string, mainAccess auth.AccessMode, ip string) *session.UserAuthSession {
	sess := &session.UserAuthSession{
		User: &session.User{
			Id:        int(userinfo.UserId),
			Name:      userinfo.Name,
			Username:  userinfo.Username,
			Extension: userinfo.Extension,
		},
		ExpiresAt:        userinfo.ExpiresAt,
		DomainId:         int(userinfo.Dc),
		Permissions:      make([]string, 0),
		License:          map[string]bool{},
		Scopes:           map[string]*session.Scope{},
		MainAccess:       mainAccess,
		MainObjClassName: mainObjClass,
		UserIp:           ip,
	}
	for _, lic := range userinfo.License {
		sess.License[lic.Id] = lic.ExpiresAt > time.Now().UnixMilli()
	}
	for _, permission := range userinfo.Permissions {
		switch auth.SuperPermission(permission.GetId()) {
		case auth.SuperCreatePermission:
			sess.SuperCreate = true
		case auth.SuperDeletePermission:
			sess.SuperDelete = true
		case auth.SuperEditPermission:
			sess.SuperEdit = true
		case auth.SuperSelectPermission:
			sess.SuperSelect = true
		}
		sess.Permissions = append(sess.Permissions, permission.GetId())
	}
	for _, scope := range userinfo.Scope {
		sess.Scopes[scope.Class] = &session.Scope{
			Id:     scope.GetId(),
			Name:   scope.GetName(),
			Abac:   scope.Abac,
			Obac:   scope.Obac,
			Rbac:   scope.Rbac,
			Class:  scope.Class,
			Access: scope.Access,
		}
	}

	for i, role := range userinfo.Roles {
		if i == 0 {
			sess.Roles = make([]*session.Role, 0)
		}
		sess.Roles = append(sess.Roles, &session.Role{
			Id:   role.GetId(),
			Name: role.GetName(),
		})
	}
	return sess
}

func getClientIp(ctx context.Context) string {
	v := ctx.Value("grpc_ctx")
	info, ok := v.(metadata.MD)
	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}
	if !ok {
		return ""
	}
	ip := strings.Join(info.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(info.Get("x-forwarded-for"), ",")
	}

	return ip
}
