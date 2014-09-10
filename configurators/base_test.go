package configurators_test

import (
	"testing"
)

import (
	"github.com/chuckpreslar/tranq/configurators"
	"github.com/chuckpreslar/tranq/serializers"
	"github.com/stretchr/testify/assert"
)

func TestNewSerializer(t *testing.T) {
	var (
		typ    = "type"
		attr   = "attribute"
		href   = "href"
		config = configurators.Base{
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
		serializer = config.NewSerializer()
	)

	assert.NotNil(t, serializer, "failed to return instance of serializers.Serializer interface")
	assert.Implements(t, (*serializers.Serializer)(nil), serializer, "failed to return instance of serializers.Serializer interface")

	var s = []string{
		config.ReservedStrings.ID,
		config.ReservedStrings.IDs,
		config.ReservedStrings.Links,
		config.ReservedStrings.Linked,
		config.ReservedStrings.Meta,
		config.ReservedStrings.Data,
		config.ReservedStrings.Type,
		config.ReservedStrings.Href,
	}

	for i := 0; i < len(s); i++ {
		assert.Equal(t, attr, s[i], "failed to map ReservedStrings with supplied AttributeNameFormatter")
	}
}

func TestFormatAttributeName(t *testing.T) {
	var (
		attr   = "attribute"
		test   = "test"
		config = configurators.Base{
			AttributeNameFormatter: serializers.NamingFormatterFunc(func(s string) string {
				return attr
			}),
		}
	)

	assert.Equal(t, attr, config.FormatAttributeName(test), "failed to format string with supplied AttributeNameFormatter")
	config.AttributeNameFormatter = nil
	assert.Equal(t, test, config.FormatAttributeName(test), "failed to return default value when no AttributeNameFormatter supplied")
}

func TestFormatTypeName(t *testing.T) {
	var (
		typ    = "type"
		test   = "test"
		config = configurators.Base{
			TypeNameFormatter: serializers.NamingFormatterFunc(func(s string) string {
				return typ
			}),
		}
	)

	assert.Equal(t, typ, config.FormatTypeName(test), "failed to format string with supplied TypeNameFormatter")
	config.TypeNameFormatter = nil
	assert.Equal(t, test, config.FormatTypeName(test), "failed to return default value when no TypeNameFormatter supplied")
}
