package tranq_test

import (
	"reflect"
	"strings"
	"testing"
)

import (
	"github.com/chuckpreslar/tranq"
	"github.com/stretchr/testify/assert"
)

func TestEmebeddedStrategyImplementation(t *testing.T) {
	assert.Implements(t, (*tranq.Strategy)(nil), &tranq.EmbeddedStrategy{}, "tranq.EmbeddedStrategy failed to implement tranq.Strategy")
}

func TestEmebeddedStrategyGetMaxMapDepth(t *testing.T) {
	var (
		maxMapDepth = 1
		strat       = &tranq.EmbeddedStrategy{MaxMapDepth: maxMapDepth}
	)

	assert.Equal(t, maxMapDepth, strat.GetMaxMapDepth(), "failed to return expected value from tranq.EmbeddedStrategy `GetMaxMapDepth`")
}

func TestEmebeddedStrategyGetMaxLinkDepth(t *testing.T) {
	var (
		maxLinkDepth = 1
		strat        = &tranq.EmbeddedStrategy{MaxLinkDepth: maxLinkDepth}
	)

	assert.Equal(t, maxLinkDepth, strat.GetMaxLinkDepth(), "failed to return expected value from tranq.EmbeddedStrategy `GetMaxLinkDepth`")
}

func TestEmebeddedStrategyGetTopLevelNamespace(t *testing.T) {
	var (
		namespace = "test"
		strat     = &tranq.EmbeddedStrategy{}
		expected  = "expected"
	)

	assert.Equal(t, expected, strat.GetTopLevelNamespace(expected), "failed to return expected value from tranq.EmbeddedStrategy `GetTopLevelNamespace` with no default set")
	strat.TopLevelNamespace = namespace
	assert.Equal(t, namespace, strat.GetTopLevelNamespace(expected), "failed to return expected value from tranq.EmbeddedStrategy `GetTopLevelNamespace` with default set")
}

func TestEmbeddedStrategyFormatAttributeName(t *testing.T) {
	var (
		attrName  = "test"
		strat     = &tranq.EmbeddedStrategy{}
		formatter = tranq.NamingFormatter(func(s string) string { return strings.ToUpper(s) })
	)

	assert.Equal(t, attrName, strat.FormatAttributeName(attrName), "failed to return expected value from tranq.EmbeddedStrategy `FormatAttributeName` with no NamingFormatter")
	strat.AttributeNameFormatter = formatter
	assert.Equal(t, formatter(attrName), strat.FormatAttributeName(attrName), "failed to return expected value from tranq.EmbeddedStrategy `FormatAttributeName` with NamingFormatter")
}

func TestEmbeddedStrategyFormatTypeName(t *testing.T) {
	var (
		attrName  = "test"
		strat     = &tranq.EmbeddedStrategy{}
		formatter = tranq.NamingFormatter(func(s string) string { return strings.ToUpper(s) })
	)

	assert.Equal(t, attrName, strat.FormatTypeName(attrName), "failed to return expected value from tranq.EmbeddedStrategy `FormatTypeName` with no NamingFormatter")
	strat.TypeNameFormatter = formatter
	assert.Equal(t, formatter(attrName), strat.FormatTypeName(attrName), "failed to return expected value from tranq.EmbeddedStrategy `FormatTypeName` with NamingFormatter")
}

func TestEmbeddedStrategyShouldSkipStructField(t *testing.T) {
	var (
		structField = reflect.StructField{
			Tag: `tranq-ignore:"true"`,
		}
		strat = &tranq.EmbeddedStrategy{}
	)

	assert.True(t, strat.ShouldSkipStructField(structField), "failed to return expected value from tranq.EmbeddedStrategy `ShouldSkipStructField`")
}

func TestEmbeddedStrategyShouldLinkStructField(t *testing.T) {
	var (
		structField = reflect.StructField{
			Tag: `tranq-link:"true"`,
		}
		strat = &tranq.EmbeddedStrategy{}
	)

	assert.True(t, strat.ShouldLinkStructField(structField), "failed to return expected value from tranq.EmbeddedStrategy `ShouldLinkStructField`")
}
