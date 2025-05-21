package model

import (
	"strings"
	"time"
)

type Session struct {
	User        *AuthorizedUser
	Permissions []*Permission
	Scope       []*Scope
	Roles       []*Role
	DomainId    int64
	ExpiresAt   int64
}

func (s *Session) HasScope(scopeName string) bool {
	for _, scope := range s.Scope {
		if scope.Name == scopeName {
			return true
		}
	}
	return false
}

func (s *Session) GetScope(scopeName string) *Scope {
	for _, scope := range s.Scope {
		if scope.Class == scopeName {
			return scope
		}
	}
	return nil
}

func (s *Session) GetUserId() int64 {
	if s.User == nil {
		return 0
	}
	return s.User.Id
}

func (s *Session) GetUser() *AuthorizedUser {
	if s.User == nil {
		return nil
	}
	clone := *s.User
	return &clone
}

func (s *Session) GetDomainId() int64 {
	return s.DomainId
}

func (s *Session) GetAclRoles() []int64 {
	roles := []int64{s.GetUserId()}
	for _, role := range s.Roles {
		roles = append(
			roles,
			role.Id,
		)
	}
	return roles
}

func (s *Session) IsExpired() bool {
	return time.Now().Unix() > s.ExpiresAt
}

func (s *Session) HasPermission(permissionName string) bool {
	for _, permission := range s.Permissions {
		if permission.Id == permissionName {
			return true
		}
	}
	return false
}

func (s *Session) HasObacAccess(scopeName string, accessType AccessMode) bool {
	scope := s.GetScope(scopeName)
	if scope == nil {
		return false
	}

	var (
		bypass, require string
	)

	switch accessType {
	case Delete, Read | Delete:
		require, bypass = "d", "delete"
	case Edit, Read | Edit:
		require, bypass = "w", "write"
	case Read, NONE:
		require, bypass = "r", "read"
	case Add, Read | Add:
		require, bypass = "x", "add"
	}
	if bypass != "" && s.HasPermission(bypass) {
		return true
	}
	for i := len(require) - 1; i >= 0; i-- {
		mode := require[i]
		if strings.IndexByte(scope.Access, mode) < 0 {
			return false
		}
	}

	return true
}

func (s *Session) UseRbacAccess(scopeName string, accessType AccessMode) bool {
	scope := s.GetScope(scopeName)
	if scope == nil || !scope.Rbac {
		return false
	}

	var (
		bypass string
	)

	switch accessType {
	case Delete, Read | Delete:
		bypass = "delete"
	case Edit, Read | Edit:
		bypass = "write"
	case Read, NONE:
		bypass = "read"
	case Add, Read | Add:
		bypass = "add"
	}
	if bypass != "" && s.HasPermission(bypass) {
		return false
	}

	return true
}
