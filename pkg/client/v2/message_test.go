package client

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRecordMarshaling(t *testing.T) {
	type args struct {
		userId     int64
		userIp     string
		recordId   int64
		recordBody any
	}

	type record struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test record marshal",
			args: args{
				userId:   0,
				userIp:   "",
				recordId: 1,
				recordBody: &record{
					Id:   1,
					Name: "name",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUpdateMessage(tt.args.userId, tt.args.userIp, tt.args.recordId, tt.args.recordBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUpdateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			body, err := json.Marshal(tt.args.recordBody)
			if err != nil {
				t.Errorf("NewUpdateMessage() error = %v", err)
				return
			}
			if got.Records[0] == nil {
				t.Errorf("NewUpdateMessage() record was not processed")
				return
			}
			if !reflect.DeepEqual(got.Records[0].NewState, body) {
				t.Errorf("bytes are not equal, want - %s, got - %s", body, got.Records[0].NewState)
				return
			}
		})
	}
}

func TestMessageMarshaling(t *testing.T) {
	type args struct {
		userId     int64
		userIp     string
		recordId   int64
		recordBody any
	}

	type record struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test record marshal",
			args: args{
				userId:   10454,
				userIp:   "192.168.0.1",
				recordId: 1,
				recordBody: &record{
					Id:   1,
					Name: "name",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUpdateMessage(tt.args.userId, tt.args.userIp, tt.args.recordId, tt.args.recordBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUpdateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			body, err := json.Marshal(got)
			if err != nil {
				t.Errorf("NewUpdateMessage() error = %v", err)
				return
			}
			t.Log(string(body))
		})
	}
}
