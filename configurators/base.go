package configurators

import (
	"github.com/chuckpreslar/tranq/serializers"
)

const (
	// JSONID ...
	JSONID = "id"
	// JSONIDs ...
	JSONIDs = "ids"
	// JSONLinks ...
	JSONLinks = "links"
	// JSONLinked ...
	JSONLinked = "linked"
	// JSONMeta ...
	JSONMeta = "meta"
	// JSONData ...
	JSONData = "data"
	// JSONHref ...
	JSONHref = "href"
	// JSONType ...
	JSONType = "type"
)

// Base ...
type Base struct {
	// TypeNameFormatter ...
	TypeNameFormatter serializers.NamingFormatter
	// AttributeNameFormatter ...
	AttributeNameFormatter serializers.NamingFormatter
	// HrefFormatter ...
	HrefFormatter serializers.HrefFormatter
	// reservedWords
	reservedWords map[string]string
}

// NewSerializer ...
func (b Base) NewSerializer() serializers.Serializer {
	if nil == b.reservedWords {
		b.reservedWords = make(map[string]string)

		b.reservedWords[JSONID] = b.FormatAttributeName(JSONID)
		b.reservedWords[JSONIDs] = b.FormatAttributeName(JSONIDs)
		b.reservedWords[JSONLinks] = b.FormatAttributeName(JSONLinks)
		b.reservedWords[JSONLinked] = b.FormatAttributeName(JSONLinked)
		b.reservedWords[JSONMeta] = b.FormatAttributeName(JSONMeta)
		b.reservedWords[JSONData] = b.FormatAttributeName(JSONData)
		b.reservedWords[JSONHref] = b.FormatAttributeName(JSONHref)
		b.reservedWords[JSONType] = b.FormatAttributeName(JSONType)
	}

	return &serializers.Base{
		TypeNameFormatter:      b.TypeNameFormatter,
		AttributeNameFormatter: b.AttributeNameFormatter,
		HrefFormatter:          b.HrefFormatter,
		ReservedWords:          b.reservedWords,
	}
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
