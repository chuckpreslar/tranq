package tranq_test

import (
	"testing"
)

import (
	"github.com/chuckpreslar/tranq"
	"github.com/chuckpreslar/tranq/configurators"
	"github.com/chuckpreslar/tranq/serializers"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		typ    = "type"
		attr   = "attribute"
		href   = "href"
		config = &configurators.Base{
			TypeNameFormatter: serializers.NamingFormatterFunc(func(s string) string {
				return typ
			}),
			AttributeNameFormatter: serializers.NamingFormatterFunc(func(s string) string {
				return attr
			}),
			HrefFormatter: serializers.HrefFormatterFunc(func(h, o, c string, i []interface{}) string {
				return href
			}),
		}
		serializer = tranq.New(config)
	)

	assert.NotNil(t, serializer, "failed to return instance of tranq.Tranq")
	assert.IsType(t, serializer, &tranq.Tranq{}, "failed to return instance of tranq.Tranq")
}

func TestSerialize(t *testing.T) {
	type TStruct struct {
		Test string
	}

	var (
		test   = "test"
		typ    = "type"
		attr   = "attribute"
		href   = "href"
		config = &configurators.Base{
			TypeNameFormatter: serializers.NamingFormatterFunc(func(s string) string {
				return typ
			}),
			AttributeNameFormatter: serializers.NamingFormatterFunc(func(s string) string {
				return attr
			}),
			HrefFormatter: serializers.HrefFormatterFunc(func(h, o, c string, i []interface{}) string {
				return href
			}),
		}
		serializer  = tranq.New(config)
		result, err = serializer.Serialize(TStruct{test})
	)

	assert.Nil(t, err, "tranq.Tranq's `Serialize` method returned an unexpected error, %s", err)
	assert.NotNil(t, result[typ], "failed to establish root level namespace for value provided to tranq.Tranq's `Serialize` method")

	result = result[typ].(map[string]interface{})
	assert.Equal(t, test, result[attr], "failed to estabish attribute returned from AttributeNameFormatter provided by configurators.Base")
}
