package tranq

import (
	"fmt"
	"strings"
)

// URLFormatter allows for formatting of linked resources
// `href` attribute.
type URLFormatter func(h, o, l string, i interface{}) string

// CommaSeparatedURLFormatter is a URLFormatter that combines
// comma separated linked resource IDs with the original `href`
// struct tag.
func CommaSeparatedURLFormatter(h, o, l string, i interface{}) string {
	h = strings.TrimRight(h, "/")

	var (
		ok    bool
		ids   []interface{}
		index int
		id    interface{}
	)

	if ids, ok = i.([]interface{}); !ok {
		return fmt.Sprintf("%s/%v", h, i)
	}

	var joined string

	for index, id = range ids {
		if index < len(ids)-1 {
			joined = fmt.Sprintf("%s%v,", joined, id)
		} else {
			joined = fmt.Sprintf("%s%v", joined, id)
		}
	}

	return fmt.Sprintf("%s/%s", h, joined)
}

// TypeBasedURLFormatter is a URLFormatter that combines
// the base type name with the linked resource name, appending
// it to the end of the original `href` struct tag.
// Ex.
//      h => "/api/comments"
//      o => "posts"
//      l => "comments"
//
//        => "/api/comments/{posts.comments}"
func TypeBasedURLFormatter(h, o, l string, i interface{}) string {
	h = strings.TrimRight(h, "/")

	return fmt.Sprintf("%s/{%s.%s}", h, o, l)
}
