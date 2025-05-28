package client

import (
	"context"
	cache "github.com/hashicorp/golang-lru/v2/expirable"
	amqp "github.com/rabbitmq/amqp091-go"
	proto "github.com/webitel/logger/api/logger"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestLoggerClient_GetObjectedLogger(t *testing.T) {
	type fields struct {
		grpcConnection   *grpc.ClientConn
		grpcClient       proto.ConfigServiceClient
		memoryCache      *cache.LRU[string, bool]
		cacheTimeToLive  time.Duration
		rabbitConnection *amqp.Connection
		channel          *amqp.Channel
	}
	type args struct {
		objclass string
	}
	objClass := "chats"
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "check same pointers and set objClass",
			fields: fields{},
			args:   args{objclass: objClass},
			want:   objClass,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoggerClient{
				grpcConnection:   tt.fields.grpcConnection,
				grpcClient:       tt.fields.grpcClient,
				memoryCache:      tt.fields.memoryCache,
				cacheTimeToLive:  tt.fields.cacheTimeToLive,
				rabbitConnection: tt.fields.rabbitConnection,
				channel:          tt.fields.channel,
			}
			if got := l.GetObjectedLogger(tt.args.objclass); !(got.parent == l || got.GetObjClass() == objClass) {
				t.Errorf("GetObjectedLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatKey(t *testing.T) {
	type args struct {
		domainId int64
		objClass string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Check formatting",
			args: args{
				domainId: 1,
				objClass: "chats",
			},
			want: "logger.1.chats",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatKey(tt.args.domainId, tt.args.objClass); got != tt.want {
				t.Errorf("formatKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLoggerClient(t *testing.T) {
	got, err := NewLoggerClient(WithGrpcConsulAddress(""), WithAmqpConnectionString(""))
	if err != nil {
		t.Error()
		return
	}

	obj := got.GetObjectedLogger("chats")
	mess, err := NewDeleteMessage(3, "", 1)
	if err != nil {
		t.Error()
		return
	}
	_, err = obj.SendContext(context.Background(), 1, mess)
	if err != nil {
		t.Error()
		return
	}
}
