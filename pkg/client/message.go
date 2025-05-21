package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type Message struct {
	records        []*InputRecord
	Records        []*Record `json:"records,omitempty"`
	RequiredFields `json:"requiredFields"`
	client         RabbitClient
}

type InputRecord struct {
	Id     int64
	Object any
}

func (r *InputRecord) TransformToOutput() (*Record, error) {
	bytes, err := json.Marshal(r.Object)
	if err != nil {
		return nil, fmt.Errorf("error marshalling record: %s", err.Error())
	}
	return &Record{
		Id:       r.Id,
		NewState: bytes,
	}, nil
}

func (c *Message) Many(records []*InputRecord) *Message {
	if len(records) == 0 {
		return c
	}
	c.records = records
	return c
}

func (c *Message) checkRecordsValidity() error {
	if c.Records == nil {
		return fmt.Errorf("logger: no records data in message")
	}
	var canNil bool
	switch c.Action {
	case CREATE_ACTION.String(), UPDATE_ACTION.String():
		canNil = false
	case DELETE_ACTION.String():
		canNil = true
	}
	if !canNil {
		for _, record := range c.Records {
			if len(record.NewState) == 0 {
				return fmt.Errorf("logger: record has no data ( id: %d )", record.Id)
			}
		}
	}

	return nil
}

func (c *Message) One(record *InputRecord) *Message {
	if record == nil {
		return c
	}
	c.records = append(c.records, record)
	return c
}

func (c *Message) SendContext(ctx context.Context) error {
	//if err := c.checkRecordsValidity(); err != nil {
	//	return err
	//}
	for _, record := range c.records {
		res, err := record.TransformToOutput()
		if err != nil {
			return err
		}
		c.Records = append(c.Records, res)
	}
	err := c.client.SendContext(ctx, c)
	if err != nil {
		return err
	}
	return nil
}
