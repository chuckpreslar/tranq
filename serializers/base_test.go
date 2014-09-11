package serializers_test

import (
	"fmt"
	"reflect"
	"testing"
)

import (
	"github.com/chuckpreslar/tranq/serializers"
	"github.com/stretchr/testify/assert"
)

func TestUninterfaceableValueError(t *testing.T) {
	var (
		val = reflect.ValueOf(1)
		err = serializers.UninterfaceableValueError{val}
		str = fmt.Sprintf("failed to call `Interface` method for reflect.Value `%s`", val)
	)

	assert.Equal(t, str, err.Error(), "failed to return correct error message for UninterfaceableValueError")
}

func TestUnsupportedKindError(t *testing.T) {
	var (
		kind       = reflect.Invalid
		serializer = &serializers.Base{}
		err        = serializers.UnsupportedKindError{kind, serializer}
		str        = fmt.Sprintf("encountered a value with reflect.Kind of `%s` is unsupported by the serializer `%T`", kind, serializer)
	)

	assert.Equal(t, str, err.Error(), "failed to return correct error message for UnsupportedKindError")
}

func TestUnlinkedResourceError(t *testing.T) {
	var (
		val = reflect.ValueOf(1)
		err = serializers.UnlinkedResourceError{val}
		str = fmt.Sprintf("value `%s` contains a nested reflect.Struct, reflect.Slice or reflect.Array which is unlinked, this is unsupported", val)
	)

	assert.Equal(t, str, err.Error(), "failed to return correct error message for UnlinkedResourceError")
}

func TestMissingIdentifierError(t *testing.T) {
	var (
		val = reflect.ValueOf(1)
		err = serializers.MissingIdentifierError{val}
		str = fmt.Sprintf("value `%s` is missing identifier field `%s`", val, serializers.ID)
	)

	assert.Equal(t, str, err.Error(), "failed to return correct error message for MissingIdentifierError")
}

func TestHrefFormatterFuncImplementation(t *testing.T) {
	var f = serializers.HrefFormatterFunc(func(h, o, c string, i []interface{}) string { return "" })
	assert.Implements(t, (*serializers.HrefFormatter)(nil), f, "HrefFormatterFunc failed to implment HrefFormatter interface")
}

func TestNamingFormatterFuncImplementation(t *testing.T) {
	var f = serializers.NamingFormatterFunc(func(s string) string { return "" })
	assert.Implements(t, (*serializers.NamingFormatter)(nil), f, "NamingFormatterFunc failed to implment NamingFormatter interface")
}

func TestNamingFormatterFunc(t *testing.T) {
}

func TestDereference(t *testing.T) {
	var (
		a = 1
		b = &a
		c = &b
	)

	var va, ta, ka, ea = serializers.Dereference(a)
	var vb, tb, kb, eb = serializers.Dereference(b)
	var vc, tc, kc, ec = serializers.Dereference(c)

	assert.Nil(t, ea, "dereferencing of type %T resulted in an error", a)
	assert.Nil(t, eb, "dereferencing of type %T resulted in an error", b)
	assert.Nil(t, ec, "dereferencing of type %T resulted in an error", c)

	assert.Equal(t, va.Int(), a, "failed to dereference %T into value %v", a, a)
	assert.Equal(t, vb.Int(), *b, "failed to dereference %T into value %v", b, *b)
	assert.Equal(t, vc.Int(), **c, "failed to dereference %T into value %v", c, **c)

	assert.Equal(t, ka, reflect.Int, "failed to dereference %T into correct kind %v", a, reflect.Int)
	assert.Equal(t, kb, reflect.Int, "failed to dereference %T into correct kind %v", b, reflect.Int)
	assert.Equal(t, kc, reflect.Int, "failed to dereference %T into correct kind %v", c, reflect.Int)

	assert.Equal(t, ta.String(), "int", "failed to dereference %T into correct type", a)
	assert.Equal(t, tb.String(), "int", "failed to dereference %T into correct type", b)
	assert.Equal(t, tc.String(), "int", "failed to dereference %T into correct type", c)
}
func TestTypeName(t *testing.T) {
	var (
		a = 1
		b = &a
		c = &b
	)

	var na, ea = serializers.TypeName(a)
	var nb, eb = serializers.TypeName(b)
	var nc, ec = serializers.TypeName(c)

	assert.Nil(t, ea, "discovering type name of type %T resulted in an error", a)
	assert.Nil(t, eb, "discovering type name of type %T resulted in an error", b)
	assert.Nil(t, ec, "discovering type name of type %T resulted in an error", c)

	assert.Equal(t, na, "int", "failed to dereference %T into correct type", a)
	assert.Equal(t, nb, "int", "failed to dereference %T into correct type", b)
	assert.Equal(t, nc, "int", "failed to dereference %T into correct type", c)
}

var serializer = serializers.Base{}

func TestAccept(t *testing.T) {
	var result, err = serializer.Accept(1)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int"], "failed to set root level namespace")
}

func TestSerialize(t *testing.T) {}

func TestSerializeInvalid(t *testing.T) {
	var _, err = serializer.SerializeInvalid(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeInvalid")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeBool(t *testing.T) {
	var result, err = serializer.Accept(true)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["bool"], "failed to set root level namespace")
	assert.True(t, result["bool"].(bool), "failed to establish correct value for SerializeBool")
}

func TestSerializeInt(t *testing.T) {
	var value = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int"], "failed to set root level namespace")
	assert.Equal(t, result["int"].(int), value, "failed to establish correct value for SerializeInt")
}

func TestSerializeInt8(t *testing.T) {
	var value int8 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int8"], "failed to set root level namespace")
	assert.Equal(t, result["int8"].(int8), value, "failed to establish correct value for SerializeInt8")
}

func TestSerializeInt16(t *testing.T) {
	var value int16 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int16"], "failed to set root level namespace")
	assert.Equal(t, result["int16"].(int16), value, "failed to establish correct value for SerializeInt16")
}

func TestSerializeInt32(t *testing.T) {
	var value int32 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int32"], "failed to set root level namespace")
	assert.Equal(t, result["int32"].(int32), value, "failed to establish correct value for SerializeInt32")
}

func TestSerializeInt64(t *testing.T) {
	var value int64 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int64"], "failed to set root level namespace")
	assert.Equal(t, result["int64"].(int64), value, "failed to establish correct value for SerializeInt64")
}

func TestSerializeUint(t *testing.T) {
	var value uint = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["uint"], "failed to set root level namespace")
	assert.Equal(t, result["uint"].(uint), value, "failed to establish correct value for SerializeUint")
}
func TestSerializeUint8(t *testing.T) {
	var value uint8 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["uint8"], "failed to set root level namespace")
	assert.Equal(t, result["uint8"].(uint8), value, "failed to establish correct value for SerializeUint8")
}

func TestSerializeUint16(t *testing.T) {
	var value uint16 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["uint16"], "failed to set root level namespace")
	assert.Equal(t, result["uint16"].(uint16), value, "failed to establish correct value for SerializeUint16")
}

func TestSerializeUint32(t *testing.T) {
	var value uint32 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["uint32"], "failed to set root level namespace")
	assert.Equal(t, result["uint32"].(uint32), value, "failed to establish correct value for SerializeUint32")
}

func TestSerializeUint64(t *testing.T) {
	var value uint64 = 1
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["uint64"], "failed to set root level namespace")
	assert.Equal(t, result["uint64"].(uint64), value, "failed to establish correct value for SerializeUint64")
}

func TestSerializeUintptr(t *testing.T) {
	var _, err = serializer.SerializeUintptr(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeUintptr")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeFloat32(t *testing.T) {
	var value float32 = 1.0
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["float32"], "failed to set root level namespace")
	assert.Equal(t, result["float32"].(float32), value, "failed to establish correct value for SerializeFloat32")
}

func TestSerializeFloat64(t *testing.T) {
	var value float64 = 1.0
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["float64"], "failed to set root level namespace")
	assert.Equal(t, result["float64"].(float64), value, "failed to establish correct value for SerializeFloat64")
}

func TestSerializeComplex64(t *testing.T) {
	var _, err = serializer.SerializeComplex64(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeComplex64")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeComplex128(t *testing.T) {
	var _, err = serializer.SerializeComplex128(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeComplex128")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeArray(t *testing.T) {
	var value = [4]int{1, 2, 3, 4}
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int"], "failed to set root level namespace")
	assert.Equal(t, len(result["int"].([]interface{})), len(value), "failed to establish correct value for SerializeArray")

	var collection = result["int"].([]interface{})

	for i := 0; i < len(collection); i++ {
		assert.Equal(t, collection[i], value[i], "failed to establish correct value for SerializeArray")
	}
}

func TestSerializeChan(t *testing.T) {
	var _, err = serializer.SerializeChan(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeChan")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeFunc(t *testing.T) {
	var _, err = serializer.SerializeFunc(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeFunc")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeInterface(t *testing.T) {
	var _, err = serializer.SerializeInterface(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeInterface")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeMap(t *testing.T) {
	var _, err = serializer.SerializeMap(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeMap")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializePtr(t *testing.T) {
	var _, err = serializer.SerializePtr(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializePtr")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestSerializeSlice(t *testing.T) {
	var value = []int{1, 2, 3, 4}
	var result, err = serializer.Accept(value)

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["int"], "failed to set root level namespace")
	assert.Equal(t, len(result["int"].([]interface{})), len(value), "failed to establish correct value for SerializeArray")

	var collection = result["int"].([]interface{})

	for i := 0; i < len(collection); i++ {
		assert.Equal(t, collection[i], value[i], "failed to establish correct value for SerializeSlice")
	}
}

func TestSerializeString(t *testing.T) {
	var result, err = serializer.Accept("string")

	assert.Nil(t, err, "received unexpected error from Accept")
	assert.NotNil(t, result["string"], "failed to set root level namespace")
	assert.Equal(t, "string", result["string"], "failed to establish correct value for SerializeString")
}

func TestSerializeStruct(t *testing.T) {
	type Post struct {
		ID   int
		Body string
	}

	var (
		id          = 1
		body        = "Lorem ipsum..."
		result, err = serializer.Accept(Post{id, body})
	)

	assert.Nil(t, err, "received unexpected error from SerializeStruct")

	var (
		post = result["Post"].(map[string]interface{})
	)

	assert.Equal(t, id, post["ID"], "failed to set correct value for field `ID`")
	assert.Equal(t, body, post["Body"], "failed to set correct value for field `Body`")
}

func TestSerializeStructMissingLink(t *testing.T) {
	type Person struct {
		ID int
	}

	type Post struct {
		Author Person
	}

	var (
		post   = Post{Person{1}}
		_, err = serializer.Accept(post)
	)

	assert.NotNil(t, err, "failed to return error from SerializeStruct with embedded unlinked resource")
}

func TestSerializeStructWithLink(t *testing.T) {
	type Person struct {
		ID int
	}

	type Post struct {
		Author Person `tranq_link:"true"`
	}

	var (
		post   = Post{Person{1}}
		_, err = serializer.Accept(post)
	)

	assert.Nil(t, err, "returned error from SerializeStruct with embedded linked resource")
}

func TestSerializeUnsafePointer(t *testing.T) {
	var _, err = serializer.SerializeUnsafePointer(reflect.ValueOf(1))

	assert.NotNil(t, err, "failed to return serializers.UnsupportedKindError from SerializeUnsafePointer")
	assert.IsType(t, err, serializers.UnsupportedKindError{}, "error was not type of serializers.UnsupportedKindError")
}

func TestLinkStructFieldSingleResource(t *testing.T) {
	type Person struct {
		ID int
	}

	type Post struct {
		Author Person
	}

	serializer.ReservedStrings.Links = "links"
	serializer.ReservedStrings.IDs = "ids"
	serializer.ReservedStrings.ID = "id"
	serializer.ReservedStrings.Href = "href"
	serializer.ReservedStrings.Type = "type"

	var (
		mapping = make(map[string]interface{})
		author  = Person{1}
		post    = Post{author}
		vpost   = reflect.ValueOf(post)
		vauthor = reflect.ValueOf(author)
		tauthor = vauthor.Type()
		kauthor = vauthor.Kind()
		fauthor = reflect.StructField{Name: "Author", Tag: `tranq_href:"/api/people"`}
		err     = serializer.LinkStructField(mapping, vpost, vauthor, tauthor, kauthor, fauthor)
	)

	assert.Nil(t, err, "received unexpected error from LinkStructField")

	var (
		mlinks  = mapping["links"].(map[string]interface{})
		mauthor = mlinks["Author"].(map[string]interface{})
	)

	assert.NotNil(t, mauthor, "failed to establish links %T", author)

	var (
		mid   = mauthor["id"]
		mtype = mauthor["type"]
		mhref = mauthor["href"]
	)

	assert.NotNil(t, mid, "failed to establish ids %T", author)
	assert.NotNil(t, mtype, "failed to establish type %T", author)
	assert.NotNil(t, mhref, "failed to establish href %T", author)

}

func TestLinkStructFieldResourceCollection(t *testing.T) {
	type Comment struct {
		ID int
	}

	type Post struct {
		Comments []Comment
	}

	serializer.ReservedStrings.Links = "links"
	serializer.ReservedStrings.IDs = "ids"
	serializer.ReservedStrings.Href = "href"
	serializer.ReservedStrings.Type = "type"

	var (
		mapping   = make(map[string]interface{})
		comments  = []Comment{Comment{1}}
		post      = Post{comments}
		vpost     = reflect.ValueOf(post)
		vcomments = reflect.ValueOf(comments)
		tcomments = vcomments.Type()
		kcomments = vcomments.Kind()
		fcomments = reflect.StructField{Name: "Comments", Tag: `tranq_href:"/api/comments"`}
		err       = serializer.LinkStructField(mapping, vpost, vcomments, tcomments, kcomments, fcomments)
	)

	assert.Nil(t, err, "received unexpected error from LinkStructField")

	var (
		mlinks  = mapping["links"].(map[string]interface{})
		mauthor = mlinks["Comments"].(map[string]interface{})
	)

	assert.NotNil(t, mauthor, "failed to establish links %T", comments)

	var (
		mids  = mauthor["ids"]
		mtype = mauthor["type"]
		mhref = mauthor["href"]
	)

	assert.NotNil(t, mids, "failed to establish ids %T", comments)
	assert.NotNil(t, mtype, "failed to establish type %T", comments)
	assert.NotNil(t, mhref, "failed to establish href %T", comments)

}

func TestFormatAttributeName(t *testing.T) {
	var (
		attr = "attribute"
		test = "test"
	)

	serializer.AttributeNameFormatter = serializers.NamingFormatterFunc(func(s string) string {
		return attr
	})

	assert.Equal(t, serializer.FormatAttributeName(test), attr, "failed to format string with supplied AttributeNameFormatter")
	serializer.AttributeNameFormatter = nil
	assert.Equal(t, serializer.FormatAttributeName(test), test, "failed to return default value when no AttributeNameFormatter supplied")
}

func TestFormatTypeName(t *testing.T) {
	var (
		typ  = "type"
		test = "test"
	)

	serializer.TypeNameFormatter = serializers.NamingFormatterFunc(func(s string) string {
		return typ
	})

	assert.Equal(t, serializer.FormatTypeName(test), typ, "failed to format string with supplied TypeNameFormatter")
	serializer.TypeNameFormatter = nil
	assert.Equal(t, serializer.FormatTypeName(test), test, "failed to return default value when no TypeNameFormatter supplied")
}

func TestFormatHref(t *testing.T) {
	var (
		href = "href"
		test = "test"
	)

	serializer.HrefFormatter = serializers.HrefFormatterFunc(func(h, o, c string, i []interface{}) string {
		return fmt.Sprintf("%s/1", href)
	})

	assert.Equal(t, serializer.FormatHref(href, "", "", []interface{}{}), fmt.Sprintf("%s/1", href), "failed to format href with supplied HrefFormatter")
	serializer.HrefFormatter = nil
	assert.Equal(t, serializer.FormatHref(test, "", "", []interface{}{}), test, "failed to return default value when no HrefFormatter supplied")
}
