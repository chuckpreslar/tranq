package serializers

import (
	"fmt"
	"reflect"
)

// UninterfaceableValueError ...
type UninterfaceableValueError struct {
	value reflect.Value
}

// Error ...
func (u UninterfaceableValueError) Error() string {
	return fmt.Sprintf("failed to call `Interface` method for reflect.Value `%s`", u.value)
}

// UnsupportedKindError ...
type UnsupportedKindError struct {
	kind       reflect.Kind
	serializer Serializer
}

// Error ..
func (u UnsupportedKindError) Error() string {
	var serializer = reflect.TypeOf(u.serializer).Elem().Name()
	return fmt.Sprintf("the reflect.Kind of `%s` is unsupported by the serializer `%s`", u.kind, serializer)
}

// UnlinkedResourceError ...
type UnlinkedResourceError struct {
	Value reflect.Value
}

func (u UnlinkedResourceError) Error() string {
	return fmt.Sprintf("value `%s` contains a nested reflect.Struct, reflect.Slice or reflect.Array which is unlinked, this is unsupported", u.Value)
}

// MissingIdentifierError ...
type MissingIdentifierError struct {
	Value reflect.Value
}

// Error ...
func (m MissingIdentifierError) Error() string {
	return fmt.Sprintf("value `%s` is missing identifier field `ID`", m.Value)
}

// HrefFormatter ...
type HrefFormatter interface {
	// FormatHref ...
	FormatHref(h, o, c string, i []interface{}) string
}

// HrefFormatterFunc ...
type HrefFormatterFunc func(h, o, c string, i []interface{}) string

// FormatHref ...
func (f HrefFormatterFunc) FormatHref(h, o, c string, i []interface{}) string {
	return f(h, o, c, i)
}

// NamingFormatter ...
type NamingFormatter interface {
	FormatName(n string) string
}

// NamingFormatterFunc ...
type NamingFormatterFunc func(n string) string

// FormatName ...
func (f NamingFormatterFunc) FormatName(n string) string {
	return f(n)
}

// Dereference ...
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

// TypeName ...
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

// Base ...
type Base struct {
	// AttributeNameFormatter ...
	AttributeNameFormatter NamingFormatter
	// TypeNameFormatter ...
	TypeNameFormatter NamingFormatter
	// HrefFormatter ...
	HrefFormatter HrefFormatter
	// RootContext ...
	RootContext map[string]interface{}
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
}

// Accept ...
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

// Serialize ...
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

// SerializeInvalid ...
func (b *Base) SerializeInvalid(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeBool ...
func (b *Base) SerializeBool(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt ...
func (b *Base) SerializeInt(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt8 ...
func (b *Base) SerializeInt8(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt16 ...
func (b *Base) SerializeInt16(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt32 ...
func (b *Base) SerializeInt32(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeInt64 ...
func (b *Base) SerializeInt64(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint ...
func (b *Base) SerializeUint(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint8 ...
func (b *Base) SerializeUint8(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint16 ...
func (b *Base) SerializeUint16(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint32 ...
func (b *Base) SerializeUint32(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUint64 ...
func (b *Base) SerializeUint64(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeUintptr ...
func (b *Base) SerializeUintptr(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeFloat32 ...
func (b *Base) SerializeFloat32(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeFloat64 ...
func (b *Base) SerializeFloat64(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeComplex64 ...
func (b *Base) SerializeComplex64(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeComplex128 ...
func (b *Base) SerializeComplex128(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeArray ...
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

// SerializeChan ...
func (b *Base) SerializeChan(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeFunc ...
func (b *Base) SerializeFunc(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeInterface ...
func (b *Base) SerializeInterface(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeMap ...
func (b *Base) SerializeMap(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializePtr ...
func (b *Base) SerializePtr(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// SerializeSlice ...
func (b *Base) SerializeSlice(v reflect.Value) (interface{}, error) {
	return b.SerializeArray(v)
}

// SerializeString ...
func (b *Base) SerializeString(v reflect.Value) (interface{}, error) {
	return v.Interface(), nil
}

// SerializeStruct ...
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
			if "true" != field.Tag.Get("tranq_link") {
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

// SerializeUnsafePointer ...
func (b *Base) SerializeUnsafePointer(v reflect.Value) (interface{}, error) {
	return nil, UnsupportedKindError{v.Kind(), b}
}

// LinkStructField ...
func (b *Base) LinkStructField(m map[string]interface{}, p, v reflect.Value, t reflect.Type, k reflect.Kind, f reflect.StructField) error {
	var (
		links map[string]interface{}
		ok    bool
	)

	if links, ok = m[b.ReservedWords.Links].(map[string]interface{}); !ok {
		links = make(map[string]interface{})
		m[b.ReservedWords.Links] = links
	}

	var (
		attr     = b.FormatAttributeName(f.Name)
		details  = make(map[string]interface{})
		href     = f.Tag.Get("tranq_href")
		typ, err = TypeName(t)
		ids      = make([]interface{}, 0, 0)
	)

	if nil != err {
		return err
	}

	typ = b.FormatTypeName(typ)

	if k == reflect.Struct {
		var id = v.FieldByName("ID")

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

			var id = element.FieldByName("ID")

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
		details[b.ReservedWords.Href] = b.FormatHref(href, parent, typ, ids)
	}

	details[b.ReservedWords.Type] = typ
	details[b.ReservedWords.IDs] = ids

	links[attr] = details

	return nil
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

// FormatHref ...
func (b *Base) FormatHref(h, o, c string, i []interface{}) string {
	if nil == b.HrefFormatter {
		return h
	}

	return b.HrefFormatter.FormatHref(h, o, c, i)
}
