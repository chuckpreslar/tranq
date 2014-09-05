package tranq

import "reflect"

// Strategy ...
type Strategy interface {
	SetPayloadRoot(p Payload)
	GetMaxMapDepth() int
	GetMaxLinkDepth() int
	GetTopLevelNamespace(s string) string
	FormatAttributeName(s string) string
	FormatTypeName(s string) string
	ShouldSkipStructField(s reflect.StructField) bool
	ShouldLinkStructField(s reflect.StructField) bool
	LinkStructField(p Payload, l Linker) error
}
