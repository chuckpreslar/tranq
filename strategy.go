package tranq

import "reflect"

// Strategy interface is used to customize how the
// serialization compiler operates.
type Strategy interface {
	// Provides the Strategy with the root Payload map,
	// allowing for direct manipulation.
	SetPayloadRoot(p Payload)
	// Returns the depth the serialization compiler should traverse
	// into a structure, mapping resources.
	GetMaxMapDepth() int
	// Returns the depth the serialization compiler should traverse
	// into a structure, linking resources.
	GetMaxLinkDepth() int
	// Returns the name of the namespace that should be used for the
	// root JSON object.
	GetTopLevelNamespace(s string) string
	// Formats a string intended to be an attribute of a JSON object.
	FormatAttributeName(s string) string
	// Formats a string intended to be used to represent types.
	FormatTypeName(s string) string
	// Informs the serialization compiler if the reflect.StructField
	// should be skipped.
	ShouldSkipStructField(s reflect.StructField) bool
	// Informs the serialization compiler if the reflect.StructField
	// should be linked.
	ShouldLinkStructField(s reflect.StructField) bool
	// Allows the Strategy to link a resource.
	LinkStructField(p Payload, l Linker) error
}
