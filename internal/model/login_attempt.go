package model

import (
	"net"
	"time"
)

type LoginAttempt struct {
	Id         int64
	Success    bool
	AuthType   string
	UserIp     *net.IPNet
	Date       time.Time
	UserId     *int
	UserName   string
	UserAgent  string
	DomainId   *int
	DomainName string
	Details    *string
}

var (
	LoginAttemptFields = struct {
		Id         string
		Success    string
		AuthType   string
		UserIp     string
		Date       string
		User       string
		UserName   string
		UserId     string
		UserAgent  string
		DomainId   string
		DomainName string
		Details    string
	}{
		Id:         "id",
		Success:    "success",
		AuthType:   "auth_type",
		UserIp:     "user_ip",
		Date:       "date",
		User:       "user",
		UserId:     "user_id",
		UserName:   "user_name",
		UserAgent:  "user_agent",
		DomainId:   "domain_id",
		DomainName: "domain_name",
		Details:    "details",
	}
)
