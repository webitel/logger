package client

import (
	"context"
	"fmt"
	"testing"
)

var noopPub = &FprintfPublisher{}

type FprintfPublisher struct{}

func (f *FprintfPublisher) Publish(_ context.Context, routingKey string, body []byte, _ map[string]any) error {
	fmt.Printf("publish routing key: %s, body: %s\n", routingKey, body)
	return nil
}

func TestNewLoggerClient(t *testing.T) {
	type args struct {
		opts []LoggerOpts
	}
	tests := []struct {
		name    string
		args    args
		want    *Logger
		wantErr bool
	}{
		{
			name: "successful creation",
			args: args{
				opts: []LoggerOpts{
					WithPublisher(noopPub),
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "failed creation",
			args: args{
				opts: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestObjectedLogger_SendContext(t *testing.T) {
	type fields struct {
		objClass string
		parent   *Logger
	}
	type args struct {
		ctx      context.Context
		domainId int64
		message  *Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful create",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      1,
						UserIp:      "1",
						Date:        1,
						Action:      CreateAction.String(),
						OperationId: "2",
					},
					Records: []*Record{
						{
							Id:       "1",
							NewState: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "successful delete",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      1,
						UserIp:      "1",
						Date:        1,
						Action:      DeleteAction.String(),
						OperationId: "2",
					},
					Records: []*Record{
						{
							Id:       "1",
							NewState: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "successful update",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      1,
						UserIp:      "1",
						Date:        1,
						Action:      UpdateAction.String(),
						OperationId: "2",
					},
					Records: []*Record{
						{
							Id:       "1",
							NewState: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no user id",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      -100,
						UserIp:      "1",
						Date:        1,
						Action:      CreateAction.String(),
						OperationId: "2",
					},
					Records: []*Record{
						{
							Id:       "1",
							NewState: nil,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no user ip",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      -100,
						UserIp:      "",
						Date:        1,
						Action:      CreateAction.String(),
						OperationId: "2",
					},
					Records: []*Record{
						{
							Id:       "1",
							NewState: nil,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no date",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      -100,
						UserIp:      "",
						Date:        -100,
						Action:      CreateAction.String(),
						OperationId: "2",
					},
					Records: []*Record{
						{
							Id:       "1",
							NewState: nil,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no records for create action",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      -100,
						UserIp:      "",
						Date:        -100,
						Action:      CreateAction.String(),
						OperationId: "2",
					},
					Records: []*Record{},
				},
			},
			wantErr: true,
		},
		{
			name: "nil action",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message: &Message{
					RequiredFields: RequiredFields{
						UserId:      -100,
						UserIp:      "",
						Date:        -100,
						Action:      CreateAction.String(),
						OperationId: "2",
					},
					Records: nil,
				},
			},
			wantErr: true,
		},

		{
			name: "without objclass",
			fields: fields{
				objClass: "",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message:  &Message{},
			},
			wantErr: true,
		},
		{
			name: "without parent logger",
			fields: fields{
				objClass: "",
				parent:   nil,
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message:  &Message{},
			},
			wantErr: true,
		},
		{
			name: "nil message",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 1,
				message:  nil,
			},
			wantErr: true,
		},
		{
			name: "zero domain",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: 0,
				message:  nil,
			},
			wantErr: true,
		},
		{
			name: " <0 domain",
			fields: fields{
				objClass: "object",
				parent:   &Logger{publisher: noopPub},
			},
			args: args{
				ctx:      context.Background(),
				domainId: -100,
				message:  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ObjectedLogger{
				objClass: tt.fields.objClass,
				parent:   tt.fields.parent,
			}
			gotOperationId, err := l.SendContext(tt.args.ctx, tt.args.domainId, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendContext() error = %v, wantErr %v, generated operation id: %s", err, tt.wantErr, gotOperationId)
				return
			}
		})
	}
}
