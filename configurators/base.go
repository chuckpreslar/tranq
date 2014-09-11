package configurators

import (
	"sync"
)

import (
	"github.com/chuckpreslar/tranq/serializers"
)

const (
	// ID represents the JSON API reserved string "id"
	ID = "id"
	// IDs represents the JSON API reserved string "ids"
	IDs = "ids"
	// Links represents the JSON API reserved string "links"
	Links = "links"
	// Linked represents the JSON API reserved string "linked"
	Linked = "linked"
	// Meta represents the JSON API reserved string "meta"
	Meta = "meta"
	// Data represents the JSON API reserved string "data"
	Data = "data"
	// Href represents the JSON API reserved string "href"
	Href = "href"
	// Type represents the JSON API reserved string "type"
	Type = "type"
)

// Base is a type implmenting the Configurator interface.
type Base struct {
	// TypeNameFormatter is used to format names of
	// types during serialization. Types include
	// base language types as well as developer
	// defined types.
	TypeNameFormatter serializers.NamingFormatter
	// AttributeNameFormatter is used to format names of
	// attributes during serialization. Attributes
	// include JSON API reserved words and struct
	// field names.
	AttributeNameFormatter serializers.NamingFormatter
	// HrefFormatter is used to format the JSON API
	// `href` attribute value when linked resources
	// are encountered during serialization.
	HrefFormatter serializers.HrefFormatter
	// ReservedStrings is a structure containing
	// JSON API reserved words formatted with the
	// AttributeNameFormatter NamingFormatter.
	ReservedStrings struct {
		ID     string
		IDs    string
		Links  string
		Linked string
		Meta   string
		Data   string
		Type   string
		Href   string
	}
	// mutex prevents mutation of exposed fields
	// used to create instances of `serlizers.Serializer`
	mutex sync.Mutex
}

// NewSerializer implements the Configurator interface
// returning an instance of the Serializer interface
// implmented by serializers.Base.
func (b *Base) NewSerializer() serializers.Serializer {
	b.mutex.Lock()

	b.ReservedStrings.ID = b.FormatAttributeName(ID)
	b.ReservedStrings.IDs = b.FormatAttributeName(IDs)
	b.ReservedStrings.Links = b.FormatAttributeName(Links)
	b.ReservedStrings.Linked = b.FormatAttributeName(Linked)
	b.ReservedStrings.Meta = b.FormatAttributeName(Meta)
	b.ReservedStrings.Data = b.FormatAttributeName(Data)
	b.ReservedStrings.Type = b.FormatAttributeName(Type)
	b.ReservedStrings.Href = b.FormatAttributeName(Href)

	var serializer = &serializers.Base{
		TypeNameFormatter:      b.TypeNameFormatter,
		AttributeNameFormatter: b.AttributeNameFormatter,
		HrefFormatter:          b.HrefFormatter,
		ReservedStrings:        b.ReservedStrings,
	}

	b.mutex.Unlock()

	return serializer
}

// FormatAttributeName allows access to Base's
// AttributeNameFormatter NameFormatter. If no
// AttributeNameFormatter was provided, the original
// string is returned in place of a formatted one.
func (b *Base) FormatAttributeName(s string) string {
	if nil == b.AttributeNameFormatter {
		return s
	}

	return b.AttributeNameFormatter.FormatName(s)
}

// FormatTypeName allows access to Base's
// TypeNameFormatter NameFormatter. If no
// TypeNameFormatter was provided, the original
// string is returned in place of a formatted one.
func (b *Base) FormatTypeName(s string) string {
	if nil == b.TypeNameFormatter {
		return s
	}

	return b.TypeNameFormatter.FormatName(s)
}
