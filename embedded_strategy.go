package tranq

import "reflect"

// EmbeddedStrategy ...
type EmbeddedStrategy struct {
	MaxMapDepth            int
	MaxLinkDepth           int
	TopLevelNamespace      string
	AttributeNameFormatter NamingFormatter
	TypeNameFormatter      NamingFormatter

	root Payload
}

// SetPayloadRoot ...
func (e *EmbeddedStrategy) SetPayloadRoot(p Payload) {
	e.root = p
}

// GetMaxMapDepth ...
func (e *EmbeddedStrategy) GetMaxMapDepth() int {
	return e.MaxMapDepth
}

// GetMaxLinkDepth ...
func (e *EmbeddedStrategy) GetMaxLinkDepth() int {
	return e.MaxLinkDepth
}

// GetTopLevelNamespace ...
func (e *EmbeddedStrategy) GetTopLevelNamespace(s string) string {
	if 0 < len(e.TopLevelNamespace) {
		return e.TopLevelNamespace
	}

	return e.FormatTypeName(s)
}

// FormatAttributeName ...
func (e *EmbeddedStrategy) FormatAttributeName(s string) string {
	if nil == e.AttributeNameFormatter {
		return s
	}

	return e.AttributeNameFormatter(s)
}

// FormatTypeName ...
func (e *EmbeddedStrategy) FormatTypeName(s string) string {
	if nil == e.TypeNameFormatter {
		return s
	}

	return e.TypeNameFormatter(s)
}

// ShouldSkipStructField ...
func (e *EmbeddedStrategy) ShouldSkipStructField(s reflect.StructField) bool {
	return s.Tag.Get("tranq-ignore") == "true"
}

// ShouldLinkStructField ...
func (e *EmbeddedStrategy) ShouldLinkStructField(s reflect.StructField) bool {
	return s.Tag.Get("tranq-link") == "true"
}

// LinkStructField ...
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
