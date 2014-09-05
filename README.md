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

__1__) Define types, creating custom serialization strategies or use those
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
  Author   Person    `tranq-link:"true" tranq-href:"/api/v1/people"`
  Comments []Comment `tranq-link:"true" tranq-href:"/api/v1/comments"`
}

```

__2__) Create and configure a serializer.

```go
strategy := &tranq.EmbeddedStrategy {
  MaxLinkDepth: 1,
  MaxMapDepth: 1,
  TypeNameFormatter: func(s string) string {
    return inflector.UnderscoreAndPluralize(s)
  },
  AttributeNameFormatter: func(s string) string {
    return inflector.Underscore(s)
  },
}

formatter := tranq.New(strategy)
```

__3__) Win.

```go
if payload, err := formatter.CompilePayload(post); nil != err {
  panic(err)
} else if result, err := payload.Marshal(); nil != err {
  panic(err)
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
            "id":1,
            "type":"people"
          },
          "comments":{  
            "href":"/v1/comments",
            "ids":[  
              1,
              2
            ],
            "type":"comments"
          }
        }
      },
      {  
        "body":"Lorem Ipsum...",
        "id":1,
        "links":{  
          "author":{  
            "href":"/v1/people",
            "id":1,
            "type":"people"
          },
          "comments":{  
            "href":"/v1/comments",
            "ids":[1, 2],
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
