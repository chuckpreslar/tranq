package tranq

import (
	"fmt"
	"reflect"
)

// NoIDFieldError ..
type NoIDFieldError struct {
	TypeName string
}

// Error ...
func (n NoIDFieldError) Error() string {
	return fmt.Sprintf("type with name `%s` is missing ID struct field supplied by tranq.Strategy's method FormatAttributeName(\"%s\")", n.TypeName, id)
}

// UnlinkableTypeError ...
type UnlinkableTypeError struct {
	Type interface{}
}

func (u UnlinkableTypeError) Error() string {
	return fmt.Sprintf("tranq.Link recieved invalid type to link, expected []tranq.Payload or tranq.Payload, was `%t`", u.Type)
}

// Linker ...
type Linker interface {
	GetIDMap() (interface{}, error)
	GetStructFieldName() string
	GetTypeName() (string, error)
	IsCollectionLink() bool
	GetStructFieldTag(s string) string
}

// Link ...
type Link struct {
	Interface   interface{}
	Value       reflect.Value
	Type        reflect.Type
	Kind        reflect.Kind
	StructField reflect.StructField
	IDFormat    string
}

// IsCollectionLink ...
func (l Link) IsCollectionLink() bool {
	return isCollectionKind(l.Kind)
}

// GetIDMap ...
func (l Link) GetIDMap() (interface{}, error) {
	if l.IsCollectionLink() {
		var (
			ids      = make([]interface{}, 0, 0)
			elements = l.Interface.([]interface{})
		)

		for _, element := range elements {
			if id, err := l.getIDAttribute(element); nil == err {
				ids = append(ids, id)
			} else {
				return nil, err
			}
		}

		return ids, nil
	}

	return l.getIDAttribute(l.Interface)
}

func (l Link) getIDAttribute(i interface{}) (interface{}, error) {
	var (
		payload Payload
		ok      bool
		id      interface{}
	)

	if payload, ok = i.(Payload); !ok {
		panic(UnlinkableTypeError{i})
	} else if id, ok = payload[l.IDFormat]; !ok {
		var (
			name string
			err  error
		)

		if name, err = TypeName(i); nil != err {
			return nil, err
		}

		return nil, NoIDFieldError{name}
	}

	return id, nil
}

// GetStructFieldName ...
func (l Link) GetStructFieldName() string {
	return l.StructField.Name
}

// GetTypeName ...
func (l Link) GetTypeName() (string, error) {
	return TypeName(l.Type)
}

// GetStructFieldTag ...
func (l Link) GetStructFieldTag(s string) string {
	return l.StructField.Tag.Get(s)
}
