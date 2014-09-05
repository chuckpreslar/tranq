package tranq_test

import (
	"fmt"
	"reflect"
	"testing"
)

import (
	"github.com/chuckpreslar/tranq"
	"github.com/stretchr/testify/assert"
)

type TStruct struct {
	ID     int
	Value  interface{}
	Nested []TStruct `tranq-link:"true"`
}

type TStrategy struct {
	RootPayload tranq.Payload
}

func (t *TStrategy) SetPayloadRoot(p tranq.Payload)       { t.RootPayload = p }
func (t *TStrategy) GetMaxMapDepth() int                  { return 1 }
func (t *TStrategy) GetMaxLinkDepth() int                 { return 1 }
func (t *TStrategy) GetTopLevelNamespace(s string) string { return s }
func (t *TStrategy) FormatAttributeName(s string) string  { return s }
func (t *TStrategy) FormatTypeName(s string) string       { return "" }

func (t *TStrategy) ShouldSkipStructField(s reflect.StructField) bool {
	return s.Tag.Get("tranq-ignore") == "true"
}

func (t *TStrategy) ShouldLinkStructField(s reflect.StructField) bool {
	return s.Tag.Get("tranq-link") == "true"
}

func (t *TStrategy) LinkStructField(p tranq.Payload, l tranq.Linker) error {
	var (
		links tranq.Payload
		ok    bool
	)

	if links, ok = p["Links"].(tranq.Payload); !ok {
		links = make(tranq.Payload)
		p["Links"] = links
	}

	var result, _ = l.GetIDMap()

	links[l.GetStructFieldName()] = result

	return nil
}

func TestInvalidKindError(t *testing.T) {
	var (
		err      = tranq.InvalidKindError{reflect.Invalid}
		actual   = err.Error()
		expected = fmt.Sprintf("an unsupported `reflect.Kind` was encoutered, was `%s`", reflect.Invalid)
	)

	assert.Equal(t, actual, expected, "failed to provide correct error message")
}

func TestUninterfaceabledValueError(t *testing.T) {
	var (
		value    = reflect.ValueOf(0)
		err      = tranq.UninterfaceabledValueError{value}
		actual   = err.Error()
		expected = fmt.Sprintf("failed to call `Interface` method on `reflect.Value` of `%v`", value)
	)

	assert.Equal(t, actual, expected, "failed to provide correct error message")
}

func TestSimpleStructDereference(t *testing.T) {
	var _, kind, _ = tranq.Dereference(TStruct{})

	assert.Equal(t, kind, reflect.Struct, "failed to result in reflect.Kind of reflect.Struct")
}

func TestPointerStructDereference(t *testing.T) {
	var _, kind, _ = tranq.Dereference(&TStruct{})

	assert.Equal(t, kind, reflect.Struct, "failed to result in reflect.Kind of reflect.Struct")
}

func TestComplexPointerStructDereference(t *testing.T) {
	var (
		tstruct    = &TStruct{}
		_, kind, _ = tranq.Dereference(&tstruct)
	)

	assert.Equal(t, kind, reflect.Struct, "failed to result in reflect.Kind of reflect.Struct")
}

func TestSimpleCollectionDereference(t *testing.T) {
	var _, kind, _ = tranq.Dereference([]TStruct{})

	assert.Equal(t, kind, reflect.Slice, "failed to result in reflect.Kind of reflect.Slice")
}

func TestPointerCollectionDereference(t *testing.T) {
	var (
		tstructs   = []*TStruct{}
		_, kind, _ = tranq.Dereference(tstructs)
	)

	assert.Equal(t, kind, reflect.Slice, "failed to result in reflect.Kind of reflect.Slice")
}

func TestComplexPointerCollectionDereference(t *testing.T) {
	var (
		tstructs   = []*TStruct{}
		_, kind, _ = tranq.Dereference(&tstructs)
	)

	assert.Equal(t, kind, reflect.Slice, "failed to result in reflect.Kind of reflect.Slice")
}

func TestSimpleStructTypeName(t *testing.T) {
	var (
		tstruct   = TStruct{}
		actual, _ = tranq.TypeName(tstruct)
		expected  = "TStruct"
	)

	assert.Equal(t, actual, expected, "failed to return correct type name for simple struct")
}

func TestPointerStructTypeName(t *testing.T) {
	var (
		tstruct   = &TStruct{}
		actual, _ = tranq.TypeName(tstruct)
		expected  = "TStruct"
	)

	assert.Equal(t, actual, expected, "failed to return correct type name for pointer struct")
}

func TestComplexPointerStructTypeName(t *testing.T) {
	var (
		tstruct   = &TStruct{}
		actual, _ = tranq.TypeName(&tstruct)
		expected  = "TStruct"
	)

	assert.Equal(t, actual, expected, "failed to return correct type name for complex pointer struct")
}

func TestSimpleCollectionTypeName(t *testing.T) {
	var (
		tstruct   = []TStruct{}
		actual, _ = tranq.TypeName(tstruct)
		expected  = "TStruct"
	)

	assert.Equal(t, actual, expected, "failed to return correct type name for simple collection")
}

func TestPointerCollectionTypeName(t *testing.T) {
	var (
		tstruct   = []*TStruct{}
		actual, _ = tranq.TypeName(tstruct)
		expected  = "TStruct"
	)

	assert.Equal(t, actual, expected, "failed to return correct type name for pointer collection")
}

func TestComplexPointerCollectionTypeName(t *testing.T) {
	var (
		tstruct   = &[]*TStruct{}
		actual, _ = tranq.TypeName(&tstruct)
		expected  = "TStruct"
	)

	assert.Equal(t, actual, expected, "failed to return correct type name for complex pointer collection")
}

func TestNew(t *testing.T) {
	var formatter *tranq.Tranq

	assert.Implements(t, (*tranq.Strategy)(nil), &TStrategy{}, "testing tranq.Strategy failed to implement interface")
	assert.NotPanics(t, func() { formatter = tranq.New(&TStrategy{}) }, "call to `New` with valid strategy panicked")
	assert.IsType(t, &tranq.Tranq{}, formatter, "failed to return pointer to an instance of tranq.Tranq")
}

func TestTranqCompilePayload(t *testing.T) {
	assert.NotPanics(t, func() { tranq.New(&TStrategy{}).CompilePayload(TStruct{}) }, "call to `CompilePayload` resulted in a panic")
}

func TestBasicStructTranqCompileCompilePayloadResult(t *testing.T) {
	var (
		formatter = tranq.New(&TStrategy{})
		tstruct   = TStruct{ID: 1, Value: "Test"}
		root, _   = formatter.CompilePayload(tstruct)
		payload   = root["TStruct"]
		actual    = payload.(tranq.Payload)["ID"]
		expected  = tstruct.ID
	)

	assert.Equal(t, expected, actual, "failed to correctly map basic TStruct to tranq.Payload")
}

func TestPointerStructTranqCompileCompilePayloadResult(t *testing.T) {
	var (
		formatter = tranq.New(&TStrategy{})
		tstruct   = TStruct{ID: 1, Value: "Test"}
		root, _   = formatter.CompilePayload(&tstruct)
		payload   = root["TStruct"]
		actual    = payload.(tranq.Payload)["ID"]
		expected  = tstruct.ID
	)

	assert.Equal(t, expected, actual, "failed to correctly map basic TStruct to tranq.Payload")
}

func TestBasicCollectionTranqCompileCompilePayloadResult(t *testing.T) {
	var (
		formatter   = tranq.New(&TStrategy{})
		nestedInner = []TStruct{TStruct{ID: 3, Value: "Inner"}}
		nestedOuter = []TStruct{TStruct{ID: 2, Value: "Inner", Nested: nestedInner}}
		tstruct     = TStruct{ID: 1, Value: "Test", Nested: nestedOuter}
		root, _     = formatter.CompilePayload(tstruct)
		payload     = root["TStruct"]
		linksTop    = payload.(tranq.Payload)["Links"]
		nested      = linksTop.(tranq.Payload)["Nested"]
	)

	if _, ok := nested.(tranq.Payload); ok {
		assert.False(t, ok, "failed to stop at tranq.Strategy's `MaxLinkDepth()`")
	}

}
