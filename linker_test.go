package tranq_test

import (
	"reflect"
	"testing"
)

import (
	"github.com/chuckpreslar/tranq"
	"github.com/stretchr/testify/assert"
)

func TestLinkGetIDMapSingleResource(t *testing.T) {
	var (
		id      = 1
		payload = tranq.Payload{"ID": id}
		value   = reflect.ValueOf(payload)
		linker  = tranq.Link{
			Interface: payload,
			Value:     value,
			Type:      value.Type(),
			Kind:      value.Kind(),
			IDFormat:  "ID",
		}
		result, _ = linker.GetIDMap()
	)

	assert.Equal(t, id, result, "failed to return expected ID from tranq.Link's `GetIDMap`")
}

func TestLinkGetIDMapResourceCollection(t *testing.T) {
	var (
		ids     = []interface{}{1, 2}
		payload = []interface{}{tranq.Payload{"ID": ids[0]}, tranq.Payload{"ID": ids[1]}}
		value   = reflect.ValueOf(payload)
		linker  = tranq.Link{
			Interface: payload,
			Value:     value,
			Type:      value.Type(),
			Kind:      value.Kind(),
			IDFormat:  "ID",
		}
		result, _ = linker.GetIDMap()
	)

	assert.Equal(t, ids, result, "failed to return expected ID from tranq.Link's `GetIDMap`")
}

func TestLinkGetStructFieldName(t *testing.T) {
	var (
		linker = tranq.Link{
			StructField: reflect.StructField{Name: "Test"},
		}
		expected = "Test"
		actual   = linker.GetStructFieldName()
	)

	assert.Equal(t, expected, actual, "failed to return expected value from tranq.Link's `GetStructFieldName`")
}

func TestLinkGetTypeName(t *testing.T) {
	type SampleStruct struct{}

	var (
		sample = SampleStruct{}
		linker = tranq.Link{
			Type: reflect.TypeOf(sample),
		}
		expected  = "SampleStruct"
		actual, _ = linker.GetTypeName()
	)

	assert.Equal(t, expected, actual, "failed to return expected value from tranq.Link's `GetTypeName`")
}

func TestLinkIsCollectionLink(t *testing.T) {
	var (
		linker = tranq.Link{
			Kind: reflect.Slice,
		}
		expected = true
		actual   = linker.IsCollectionLink()
	)

	assert.Equal(t, expected, actual, "failed to return expected value from tranq.Link's `IsCollectionLink`")
}
func TestLinkGetStructFieldTag(t *testing.T) {
	var (
		structField = reflect.StructField{
			Tag: `test:"value"`,
		}
		linker = tranq.Link{
			StructField: structField,
		}
		expected = "value"
		actual   = linker.GetStructFieldTag("test")
	)

	assert.Equal(t, expected, actual, "failed to return expected value from tranq.Link's `GetStructFieldTag`")
}
