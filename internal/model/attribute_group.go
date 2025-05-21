package model

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
)

type AttributeGroup struct {
	attrs    []attribute.KeyValue
	rootName string
}

func NewAttributeGroup(rootName string, attrs ...attribute.KeyValue) *AttributeGroup {
	return &AttributeGroup{
		rootName: rootName,
		attrs:    attrs,
	}
}

func (a *AttributeGroup) Unparse() []attribute.KeyValue {
	var res []attribute.KeyValue
	for _, attr := range a.attrs {
		res = append(res, attribute.KeyValue{Key: attribute.Key(fmt.Sprintf("%s.%s", a.rootName, attr.Key)), Value: attr.Value})
	}
	return res
}
