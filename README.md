tranq
=====

Easily serialize and marshal responses for your self describing RESTful/JSON API
(based on [JSON API](http://jsonapi.org/)).

[![Build Status](https://drone.io/github.com/chuckpreslar/tranq/status.png)](https://drone.io/github.com/chuckpreslar/tranq/latest)

### Configuration and Extension

Due to the variety of options for linking resources specified by
[JSON API](http://jsonapi.org/), the `tranq` package allows for easy adaptation
and customization to meet the methods used by your applications through
use of Go's interfaces.

#### Configurators

`Configurators` contained in the `configurators` package serve as a layer of
abstraction between the `tranq` and the `serializers` packages. Types
implementing the `Configurator` interface are used to instantiate the `Tranq`
type and serve as factories for creating instances of the `Serializer`
interface.

#### Serializers

`Serializers` fulfill the role of accepting both native Go and user defined
types and serializing them into a Go map of the structure specified
by [JSON API](http://jsonapi.org/) for representing resources. These maps
can then be passed to the `encoding/json` package's `Marshal` method for
JSON encoding.

#### Base

Included in both the `configurators` and `serializers` package are `Base`
types. With these two starting points, users can easily compose new
types, overriding specific functionality (i.e. exactly how where a linked
resource should live, at the root level level or embedded in a document), or
provided additional functionality all together.

### Usage

__1__) Define types, create custom serialization strategies or use those
included with the package. The `tranq` package makes use of Go struct tags
to define linked resources to include with the response.

```go
type Person struct {
  Id        int
  FirstName string
  LastName  string
}

type Comment struct {
  Id     int
  Body   string
  Author Person
}

type Post struct {
  Id       int
  Body     string
  Author   Person    `tranq_link:"true" tranq_href:"/api/v1/people"`
  Comments []Comment `tranq_link:"true" tranq_href:"/api/v1/comments"`
}

```

__2__) Import, create and configure a serializer.

```go
package main

import (
  "fmt"
  "strings"
)

import (
  "github.com/chuckpreslar/inflect"
  "github.com/chuckpreslar/tranq"
  "github.com/chuckpreslar/tranq/configurators"
  "github.com/chuckpreslar/tranq/serializers"
)

func FormatTypeName(s string) string {
  return inflect.Underscore(inflect.Pluralize(s))
}

func FormatAttributeName(s string) string {
  return inflect.Underscore(s)
}

func FormatHref(href, owner, child string, ids []interface{}) string {
  var (
    str = ""
    il  = len(ids)
  )

  for i := 0; i < il; i++ {
    if i != il - 1 {
      str = fmt.Sprintf("%s%v,", str, ids[i])
			continue
		}

    str = fmt.Sprintf("%s%v", str, ids[i])
  }

  return fmt.Sprintf("%s/%s", strings.TrimRight(href, "/"), str)
}

func main() {
  configuration := new(configurators.Base)
  configuration.TypeNameFormatter = serializers.NamingFormatterFunc(FormatTypeName)
  configuration.AttributeNameFormatter = serializers.NamingFormatterFunc(FormatAttributeName)
  configuration.HrefFormatter = serializers.HrefFormatterFunc(FormatHref)
  serializer := tranq.New(configuration)
}
```

__3__) Win.

```go
// ...

if result, err := serializer.Serialize(posts); nil != err {
  return err
} else {
  return json.Marshal(result)
}

/**
  result =>
  {  
    "posts":[  
      {  
        "body":"Lorem Ipsum...",
        "id":1,
        "links":{  
          "author":{  
            "href":"/api/v1/people",
            "id": 1,
            "type":"people"
          },
          "comments":{  
            "href":"/api/v1/comments/1,2",
            "ids":[1, 2],
            "type":"comments"
          }
        }
      },
      {  
        "body":"Lorem Ipsum...",
        "id":2,
        "links":{  
          "author":{  
            "href":"/api/v1/people/2",
            "id":2,
            "type":"people"
          },
          "comments":{  
            "href":"/api/v1/comments/3,4",
            "ids":[3, 4],
            "type":"comments"
          }
        }
      }
    ]
  }
*/
```

### Documentation

View godoc or visit [godoc.org](http://godoc.org/github.com/chuckpreslar/tranq).

    $ godoc tranq

### License

> The MIT License (MIT)

> Copyright (c) 2014 Chuck Preslar

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.
