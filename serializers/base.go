package serializers

import (
	"fmt"
	"reflect"
)

const (
	// ID represents the name of a field intended
	// to be used as a resources identifier.
	ID = "ID"
)

const (
	// TranqLink represents the struct tag expected
	// for linking nested resources.
	TranqLink = "tranq_link"
	// TranqHref represents the struct tag that
	// contains the unformatted JSON API href
	// attribute.
	TranqHref = "tranq_href"
)

// UninterfaceableValueError occurs when a reflect.Value
// cannot have its `Interface` method called during
// serialization.
type UninterfaceableValueError struct {
	value reflect.Value
}

// Error implements the `error` interface for the
// UninterfaceableValueError type.
func (u UninterfaceableValueError) Error() string {
	return fmt.Sprintf("failed to call `Interface` method for reflect.Value `%s`", u.value)
}

// UnsupportedKindError occurs when a Serializer encounters
// a reflect.Kind it does not have the ability to interact
// with.
type UnsupportedKindError struct {
	Kind       reflect.Kind
	Serializer Serializer
}

// Error implements the `error` interface for the
// UnsupportedKindError type.
func (u UnsupportedKindError) Error() string {
	var serializer = reflect.TypeOf(u.Serializer).Elem().Name()
	return fmt.Sprintf("the reflect.Kind of `%s` is unsupported by the serializer `%s`", u.Kind, serializer)
}

// UnlinkedResourceError occurs when a resource that
// should be linked is encounted but is missing the
// required struct tag to mark it as so.
type UnlinkedResourceError struct {
	Value reflect.Value
}

// Error implements the `error` interface for the
// UnlinkedResourceError type.
func (u UnlinkedResourceError) Error() string {
	return fmt.Sprintf("value `%s` contains a nested reflect.Struct, reflect.Slice or reflect.Array which is unlinked, this is unsupported", u.Value)
}

// MissingIdentifierError occurs when a resource
// is flagged for linking via a struct tag but
// has no field named by the constant string
// contained in ID.
type MissingIdentifierError struct {
	Value reflect.Value
}

// Error implements the `error` interface for the
// MissingIdentifierError type.
func (m MissingIdentifierError) Error() string {
	return fmt.Sprintf("value `%s` is missing identifier field `%s`", m.Value, ID)
}

// HrefFormatter provides an interface for formatting
// JSON API linked resources `href` attribute.
type HrefFormatter interface {
	// FormatHref ...
	FormatHref(href, owner, child string, ids []interface{}) string
}

// HrefFormatterFunc is an adapter to allow the use of
// ordinary functions as HrefFormatters. If f is a function
// with the appropriate signature, HrefFormatterFunc(f)
// is a HrefFormatter object that calls f.
type HrefFormatterFunc func(h, o, c string, i []interface{}) string

// FormatHref calls f(h,o,c,i)
func (f HrefFormatterFunc) FormatHref(h, o, c string, i []interface{}) string {
	return f(h, o, c, i)
}

// NamingFormatter provides an interface for formatting
// names of attributes and types stored in a JSON API
// response.
type NamingFormatter interface {
	FormatName(name string) string
}

// NamingFormatterFunc is an adapter to allow the use of
// ordinary functions as NamingFormatters. If f is a function
// with the appropriate signature, NamingFormatterFunc(f)
// is a NamingFormatter object that calls f.
type NamingFormatterFunc func(n string) string

// FormatName calls f(n)
func (f NamingFormatterFunc) FormatName(n string) string {
	return f(n)
}

// Dereference attempts to dereference argument `i` from
// a pointer or interface to a base type, returning its
// reflect.Value, reflect.Type and reflect.Kind. If
// a reflect.Value cannot have its `Interface` method
// called without panicking, an UninterfaceableValueError
// is returned.
func Dereference(i interface{}) (reflect.Value, reflect.Type, reflect.Kind, error) {
	var (
		v reflect.Value
		k reflect.Kind
		t reflect.Type
	)

	v = reflect.ValueOf(i)
	k = v.Kind()

	if k == reflect.Ptr || k == reflect.Interface {
		v = v.Elem()

		if v.CanInterface() {
			return Dereference(v.Interface())
		}

		return reflect.Value{}, nil, reflect.Invalid, UninterfaceableValueError{v}
	}

	t = v.Type()

	return v, t, k, nil
}

// TypeName attempts to resolved the name of a type,
// both native and user defined. If argument `i`
// cannot be successfully passed to Dereference,
// an UninterfaceableValueError is returned.
func TypeName(i interface{}) (string, error) {
	var (
		o bool
		t reflect.Type
		k reflect.Kind
		e error
	)

	if t, o = i.(reflect.Type); !o {
		_, t, k, e = Dereference(i)
		if nil != e {
			return "", e
		}
	} else {
		k = t.Kind()
	}

	if k == reflect.Array || k == reflect.Slice || k == reflect.Ptr {
		return TypeName(t.Elem())
	}

	return t.Name(), nil
}

// Base is a type implementing the Serializer interface.
type Base struct {
	// TypeNameFormatter is used to format names of
	// types during serialization. Types include
	// base language types as well as developer
	// defined types.
	TypeNameFormatter NamingFormatter
	// AttributeNameFormatter is used to format names of
	// attributes during serialization. Attributes
	// include JSON API reserved words and struct
	// field names.
	AttributeNameFormatter NamingFormatter
	// HrefFormatter is used to format the JSON API
	// `href` attribute value when linked resources
	// are encountered during serialization.
	HrefFormatter HrefFormatter
	// RootContext is the base map[string]interface{}
	// created to contain the serialized JSON API
	// response.
	RootContext map[string]interface{}
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
}

// Accept implements the `Accept` method required
// by the Serializer interface.
func (b *Base) Accept(i interface{}) (map[string]interface{}, error) {
	var (
		err       error
		namespace string
		mapping   map[string]interface{}
	)

	defer func() {
		if temp := recover(); nil != temp {
			if _, ok := temp.(error); ok {
				err = temp.(error)
			} else {
				err = fmt.Errorf("%s", temp)
			}

			mapping = nil
		}
	}()

	if namespace, err = TypeName(i); nil != err {
		return nil, err
	}

	namespace = b.FormatTypeName(namespace)
	mapping = make(map[string]interface{})
	b.RootContext = mapping

	mapping[namespace], err = b.Serialize(i)

	return mapping, err
}

// Serialize allows for the recursive serialization
// of base and user defeined types.
func (b *Base) Serialize(i interface{}) (interface{}, error) {
	var (
		result interface{}
		value  reflect.Value
		kind   reflect.Kind
		err    error
	)

	if value, _, kind, err = Dereference(i); nil != err {
		return nil, err
	} else if !value.CanInterface() {
		return nil, UninterfaceableValueError{value}
	}

	switch kind {
	case reflect.Invalid:
		result, err = b.SerializeInvalid(value)
	case reflect.Bool:
		result, err = b.SerializeBool(value)
	case reflect.Int:
		result, err = b.SerializeInt(value)
	case reflect.Int8:
		result, err = b.SerializeInt8(value)
	case reflect.Int16:
		result, err = b.SerializeInt16(value)
	case reflect.Int32:
		result, err = b.SerializeInt32(value)
	case reflect.Int64:
		result, err = b.SerializeInt64(value)
	case reflect.Uint:
		result, err = b.SerializeUint(value)
	case reflect.Uint8:
		result, err = b.SerializeUint8(value)
	case reflect.Uint16:
		result, err = b.SerializeUint16(value)
	case reflect.Uint32:
		result, err = b.SerializeUint32(value)
	case reflect.Uint64:
		result, err = b.SerializeUint64(value)
	case reflect.Uintptr:
		result, err = b.SerializeUintptr(value)
	case reflect.Float32:
		result, err = b.SerializeFloat32(value)
	case reflect.Float64:
		result, err = b.SerializeFloat64(value)
	case reflect.Complex64:
		result, err = b.SerializeComplex64(value)
	case reflect.Complex128:
		result, err = b.SerializeComplex128(value)
	case reflect.Array:
		result, err = b.SerializeArray(value)
	case reflect.Chan:
		result, err = b.SerializeChan(value)
	case reflect.Func:
		result, err = b.SerializeFunc(value)
	case reflect.Interface:
		result, err = b.SerializeInterface(value)
	case reflect.Map:
		result, err = b.SerializeMap(value)
	case reflect.Ptr:
		result, err = b.SerializePtr(value)
	case reflect.Slice:
		result, err = b.SerializeSlice(value)
	case reflect.String:
		result, err = b.SerializeString(value)
	case reflect.Struct:
		result, err = b.SerializeStruct(value)
	case reflect.UnsafePointer:
		result, err = b.SerializeUnsafePointer(value)
	}

	return result, err
}

// SerializeInvalid attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Invalid.
func (b *Base) SerializeInvalid(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeBool attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Bool.
func (b *Base) SerializeBool(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Int.
func (b *Base) SerializeInt(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt8 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Int8.
func (b *Base) SerializeInt8(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt16 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Int16.
func (b *Base) SerializeInt16(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt32 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Int32.
func (b *Base) SerializeInt32(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt64 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Int64.
func (b *Base) SerializeInt64(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Uint.
func (b *Base) SerializeUint(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint8 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Uint8.
func (b *Base) SerializeUint8(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint16 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Uint16.
func (b *Base) SerializeUint16(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint32 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Uint32.
func (b *Base) SerializeUint32(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint64 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Uint64.
func (b *Base) SerializeUint64(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUintptr attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Uintptr.
func (b *Base) SerializeUintptr(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeFloat32 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Float32.
func (b *Base) SerializeFloat32(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeFloat64 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Float64.
func (b *Base) SerializeFloat64(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeComplex64 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Complex64.
func (b *Base) SerializeComplex64(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeComplex128 attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Complex128.
func (b *Base) SerializeComplex128(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeArray attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Array.
func (b *Base) SerializeArray(v reflect.Value) (interface{}, error) {
	var collection = make([]interface{}, 0, 0)

	for i := 0; i < v.Len(); i++ {
		var element = v.Index(i)

		if !element.CanInterface() {
			panic(UninterfaceableValueError{element})
		}

		var result, err = b.Serialize(element.Interface())

		if nil != err {
			return nil, err
		}

		collection = append(collection, result)
	}

	return collection, nil
}

// SerializeChan attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Chan.
func (b *Base) SerializeChan(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeFunc attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Func.
func (b *Base) SerializeFunc(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeInterface attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Interface.
func (b *Base) SerializeInterface(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeMap attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Map.
func (b *Base) SerializeMap(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializePtr attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Ptr.
func (b *Base) SerializePtr(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeSlice attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Slice.
func (b *Base) SerializeSlice(v reflect.Value) (interface{}, error) {
	return b.SerializeArray(v)
}

// SerializeString attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.String.
func (b *Base) SerializeString(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeStruct attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.Struct.
func (b *Base) SerializeStruct(v reflect.Value) (interface{}, error) {
	var (
		mapping = make(map[string]interface{})
		t       = v.Type()
	)

	for i := 0; i < v.NumField(); i++ {
		var temp = v.Field(i)

		if !temp.CanInterface() {
			return nil, UninterfaceableValueError{temp}
		}

		var val, typ, kind, err = Dereference(temp.Interface())

		if nil != err {
			return nil, err
		}

		var (
			field = t.Field(i)
			attr  = b.FormatAttributeName(field.Name)
		)

		if kind == reflect.Struct || kind == reflect.Array || kind == reflect.Slice {
			if "true" != field.Tag.Get(TranqLink) {
				return nil, UnlinkedResourceError{v}
			} else if err = b.LinkStructField(mapping, v, val, typ, kind, field); nil != err {
				return nil, err
			}

		} else if mapping[attr], err = b.Serialize(val.Interface()); nil != err {
			return nil, err
		}
	}

	return mapping, nil
}

// SerializeUnsafePointer attempts to serialize a reflect.Value with a reflect.Kind
// of reflect.UnsafePointer.
func (b *Base) SerializeUnsafePointer(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// LinkStructField attempts to add link details for a value
// to a map[string]interface{} under the JSON API reserved
// string `links`. If LinkStructField an error detailing
// what went wrong is returned.
func (b *Base) LinkStructField(m map[string]interface{}, p, v reflect.Value, t reflect.Type, k reflect.Kind, f reflect.StructField) error {
	var (
		links map[string]interface{}
		ok    bool
	)

	if links, ok = m[b.ReservedStrings.Links].(map[string]interface{}); !ok {
		links = make(map[string]interface{})
		m[b.ReservedStrings.Links] = links
	}

	var (
		attr     = b.FormatAttributeName(f.Name)
		details  = make(map[string]interface{})
		href     = f.Tag.Get(TranqHref)
		typ, err = TypeName(t)
		ids      = make([]interface{}, 0, 0)
	)

	if nil != err {
		return err
	}

	typ = b.FormatTypeName(typ)

	if k == reflect.Struct {
		var id = v.FieldByName(ID)

		if !id.IsValid() {
			return MissingIdentifierError{v}
		}

		ids = append(ids, id.Interface())
	} else if k == reflect.Slice || k == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			var temp = v.Index(i)

			if !temp.CanInterface() {
				return UninterfaceableValueError{temp}
			}

			var element, _, _, err = Dereference(temp.Interface())

			if nil != err {
				return err
			}

			var id = element.FieldByName(ID)

			if !id.IsValid() {
				return MissingIdentifierError{v}
			}

			ids = append(ids, id.Interface())
		}
	}

	if 0 < len(href) {
		var parent string

		if parent, err = TypeName(p.Type()); nil != err {
			return err
		}

		parent = b.FormatTypeName(parent)
		details[b.ReservedStrings.Href] = b.FormatHref(href, parent, typ, ids)
	}

	details[b.ReservedStrings.Type] = typ
	details[b.ReservedStrings.IDs] = ids

	links[attr] = details

	return nil
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

// FormatHref allows access to Base's HrefFormatter
// HrefFormatter. If no FormatHref was provided, the
// original `href` string is returned in place of a
// formatted one.
func (b *Base) FormatHref(h, o, c string, i []interface{}) string {
	if nil == b.HrefFormatter {
		return h
	}

	return b.HrefFormatter.FormatHref(h, o, c, i)
}
