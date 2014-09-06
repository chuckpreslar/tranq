package tranq_test

import (
	"testing"

	"github.com/chuckpreslar/tranq"
	"github.com/stretchr/testify/assert"
)

func TestCommaSeparatedURLFormatter(t *testing.T) {
	var (
		id   = 1
		ids  = []interface{}{1, 2, 3, 4}
		href = "/api/test"
	)

	var (
		expected = "/api/test/1"
		actual   = tranq.CommaSeparatedURLFormatter(href, "", "", id)
	)

	assert.Equal(t, expected, actual, "failed to return expected result from CommaSeparatedURLFormatter with single id")

	expected = "/api/test/1,2,3,4"
	actual = tranq.CommaSeparatedURLFormatter(href, "", "", ids)

	assert.Equal(t, expected, actual, "failed to return expected result from CommaSeparatedURLFormatter with multiple ids")
}

func TestTypeBasedURLFormatter(t *testing.T) {
	var (
		href     = "/api/comments"
		base     = "posts"
		linked   = "comments"
		expected = "/api/comments/{posts.comments}"
		actual   = tranq.TypeBasedURLFormatter(href, base, linked, nil)
	)

	assert.Equal(t, expected, actual, "faled to return expected result from TypeBasedURLFormatter, expected `%s`, was `%s`", expected, actual)
}
