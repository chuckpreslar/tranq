package tranq

import "reflect"

// EmbeddedStrategy implements the Strategy interface,
// formatting JSON responses with a JSON API `links`
// attribute established within each resource.
type EmbeddedStrategy struct {
	MaxMapDepth            int
	MaxLinkDepth           int
	TopLevelNamespace      string
	AttributeNameFormatter NamingFormatter
	TypeNameFormatter      NamingFormatter
	URLFormatter           URLFormatter

	root Payload
}

// SetPayloadRoot sets the unexported root Payload field
func (e *EmbeddedStrategy) SetPayloadRoot(p Payload) {
	e.root = p
}

// GetMaxMapDepth returns the MaxMapDepth integer to inform
// the serialization compiler how far to traverse into
// the provided object, mapping nested resources.
func (e *EmbeddedStrategy) GetMaxMapDepth() int {
	return e.MaxMapDepth
}

// GetMaxLinkDepth returns the MaxLinkDepth integer to inform
// the serialization compiler how far to traverse into
// the provided object, linking nested resources.
func (e *EmbeddedStrategy) GetMaxLinkDepth() int {
	return e.MaxLinkDepth
}

// GetTopLevelNamespace returns the provided top level namespace
// for the root mapping if one is specified, otherwise it attempts
// to format the string `s` with the method `FormatTypeName`.
func (e *EmbeddedStrategy) GetTopLevelNamespace(s string) string {
	if 0 < len(e.TopLevelNamespace) {
		return e.TopLevelNamespace
	}

	return e.FormatTypeName(s)
}

// FormatAttributeName returns the result of the AttributeNameFormatter
// if one is provided, else the original string `s` is returned.
func (e *EmbeddedStrategy) FormatAttributeName(s string) string {
	if nil == e.AttributeNameFormatter {
		return s
	}

	return e.AttributeNameFormatter(s)
}

// FormatTypeName returns the result of the TypeNameFormatter
// if one is provided, else the original string `s` is returned.
func (e *EmbeddedStrategy) FormatTypeName(s string) string {
	if nil == e.TypeNameFormatter {
		return s
	}

	return e.TypeNameFormatter(s)
}

// ShouldSkipStructField informs the serialization compiler if a struct
// field should be skipped over.
func (e *EmbeddedStrategy) ShouldSkipStructField(s reflect.StructField) bool {
	return s.Tag.Get("tranq-ignore") == "true"
}

// ShouldLinkStructField informs the serialization compiler if a struct
// field should be linked.
func (e *EmbeddedStrategy) ShouldLinkStructField(s reflect.StructField) bool {
	return s.Tag.Get("tranq-link") == "true"
}

// LinkStructField links a resource using the Linker interface to the
// supplied Payload.
func (e *EmbeddedStrategy) LinkStructField(p Payload, l Linker) error {
	var (
		ok    bool
		links Payload
		ids   interface{}
	)

	var (
		linksStr = e.FormatAttributeName("links")
		typeStr  = e.FormatAttributeName("type")
		idStr    = e.FormatAttributeName("id")
		idsStr   = e.FormatAttributeName("ids")
		hrefStr  = e.FormatAttributeName("href")
	)

	if links, ok = p[linksStr].(Payload); !ok {
		links = make(Payload)
		p[linksStr] = links
	}

	var (
		typeName, err = l.GetTypeName()
		attributeName = e.FormatAttributeName(l.GetStructFieldName())
		details       = make(Payload)
	)

	if nil != err {
		return err
	}

	typeName = e.FormatTypeName(typeName)
	details[typeStr] = typeName

	if ids, err = l.GetIDMap(); nil != err {
		return err
	}

	if l.IsCollectionLink() {
		details[idsStr] = ids
	} else {
		details[idStr] = ids
	}

	if href := l.GetStructFieldTag("tranq-href"); 0 < len(href) {
		details[hrefStr] = href
	}

	links[attributeName] = details

	return nil
}
