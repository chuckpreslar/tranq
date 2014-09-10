package configurators

import (
	"sync"

	"github.com/chuckpreslar/tranq/serializers"
)

const (
	// ID ...
	ID = "id"
	// IDs ...
	IDs = "ids"
	// Links ...
	Links = "links"
	// Linked ...
	Linked = "linked"
	// Meta ...
	Meta = "meta"
	// Data ...
	Data = "data"
	// Href ...
	Href = "href"
	// Type ...
	Type = "type"
)

// Base ...
type Base struct {
	// TypeNameFormatter ...
	TypeNameFormatter serializers.NamingFormatter
	// AttributeNameFormatter ...
	AttributeNameFormatter serializers.NamingFormatter
	// HrefFormatter ...
	HrefFormatter serializers.HrefFormatter
	// ReservedWords ...
	ReservedWords struct {
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

// NewSerializer ...
func (b Base) NewSerializer() serializers.Serializer {
	b.mutex.Lock()

	b.ReservedWords.ID = b.FormatAttributeName(ID)
	b.ReservedWords.IDs = b.FormatAttributeName(IDs)
	b.ReservedWords.Links = b.FormatAttributeName(Links)
	b.ReservedWords.Linked = b.FormatAttributeName(Linked)
	b.ReservedWords.Meta = b.FormatAttributeName(Meta)
	b.ReservedWords.Data = b.FormatAttributeName(Data)
	b.ReservedWords.Type = b.FormatAttributeName(Type)
	b.ReservedWords.Href = b.FormatAttributeName(Href)

	var serializer = &serializers.Base{
		TypeNameFormatter:      b.TypeNameFormatter,
		AttributeNameFormatter: b.AttributeNameFormatter,
		HrefFormatter:          b.HrefFormatter,
		ReservedWords:          b.ReservedWords,
	}

	b.mutex.Unlock()

	return serializer
}

// FormatAttributeName ...
func (b *Base) FormatAttributeName(s string) string {
	if nil == b.AttributeNameFormatter {
		return s
	}

	return b.AttributeNameFormatter.FormatName(s)
}

// FormatTypeName ...
func (b *Base) FormatTypeName(s string) string {
	if nil == b.TypeNameFormatter {
		return s
	}

	return b.TypeNameFormatter.FormatName(s)
}
