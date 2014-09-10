tranq
=====

Easily format and marshal responses for your self describing RESTful/JSON API
(based on [JSON API](http://jsonapi.org/)).

[![Build Status](https://drone.io/github.com/chuckpreslar/tranq/status.png)](https://drone.io/github.com/chuckpreslar/tranq/latest)

### Configuration

Configuring the `tranq` package is done through the implementation of the
`tranq.Strategy` interface. This allows custom compilation of Go structs and
collections (arrays and slices) into specially formatted JSON to meet the
standards set by [JSON API](http://jsonapi.org/).

### Usage

__1__) Define types, create custom serialization strategies or use those
included with the package.

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

func main() {
  configuration := configurators.Base {
    TypeNameFormatter: serializers.NamingFormatterFunc(FormatTypeName),
    AttributeNameFormatter: serializers.NamingFormatterFunc(FormatAttributeName),
  }

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
            "href":"/v1/people",
            "ids":[1],
            "type":"people"
          },
          "comments":{  
            "href":"/v1/comments",
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
            "href":"/v1/people",
            "ids":[2],
            "type":"people"
          },
          "comments":{  
            "href":"/v1/comments",
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

> Copyright (c) 2013 Chuck Preslar

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
