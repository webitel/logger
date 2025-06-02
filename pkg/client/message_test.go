package client

import (
	"encoding/json"
	"testing"
)

func TestMessageMarshaling(t *testing.T) {
	type args struct {
		userId     int64
		userIp     string
		recordId   string
		recordBody any
	}

	type record struct {
		Id   string `json:"id"`
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
				recordId: "1",
				recordBody: &record{
					Id:   "1",
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
