# qstring

This package provides an easy way to marshal and unmarshal url query string data to and from structs.

This was originally forked from [dyninc/qstring](https://github.com/dyninc/qstring) but it seems to have become inactive
and I wanted some extra features so I am now maintaining this fork.

## Installation

```bash
$ go get github.com/Southclaws/qstring
```

## Examples

### Unmarshaling

```go
package main

import (
	"net/http"

	"github.com/Southclaws/qstring"
)

// Query is the http request query struct.
type Query struct {
	Names    []string
	Limit     int
	Page      int
}

func handler(w http.ResponseWriter, req *http.Request) {
	query := &Query{}
	err := qstring.Unmarshal(req.Url.Query(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// ... run conditional logic based on provided query parameters
}
```

The above example will unmarshal the query string from an http.Request and unmarshal it into the provided struct. This
means that a query of `?names=foo&names=bar&limit=50&page=1` would be unmarshaled into a struct similar to the
following:

```go
Query{
	Names: []string{"foo", "bar"},
	Limit: 50,
	Page: 1,
}
```

### Marshalling

`qstring` also exposes two methods of Marshaling structs _into_ Query parameters, one will Marshal the provided struct
into a raw query string, the other will Marshal a struct into a `url.Values` type. Some Examples of both follow.

### Marshal Raw Query String

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/Southclaws/qstring"
)

// Query is the http request query struct.
type Query struct {
	Names    []string
	Limit     int
	Page      int
}

func main() {
	query := &Query{
		Names: []string{"foo", "bar"},
		Limit: 50,
		Page: 1,
	}
	q, err := qstring.MarshalString(query)
	fmt.Println(q)
	// Output: names=foo&names=bar&limit=50&page=1
}
```

### Marshal url.Values

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/Southclaws/qstring"
)

// Query is the http request query struct.
type Query struct {
	Names    []string
	Limit     int
	Page      int
}

func main() {
	query := &Query{
		Names: []string{"foo", "bar"},
		Limit: 50,
		Page: 1,
	}
	q, err := qstring.Marshal(query)
	fmt.Println(q)
	// Output: map[names:[foo, bar] limit:[50] page:[1]]
}
```

### Nested

In the same spirit as other Unmarshaling libraries, `qstring` allows you to Marshal/Unmarshal nested structs

```go
package main

import (
	"net/http"

	"github.com/Southclaws/qstring"
)

// PagingParams represents common pagination information for query strings
type PagingParams struct {
	Page int
	Limit int
}

// Query is the http request query struct.
type Query struct {
	Names    []string
	PageInfo PagingParams
}
```

### Complex Structures

Again, in the spirit of other Unmarshaling libraries, `qstring` allows for some more complex types, such as pointers and
time.Time fields. A more complete example might look something like the following code snippet

```go
package main

import (
	"time"
)

// PagingParams represents common pagination information for query strings
type PagingParams struct {
	Page int	`qstring:"page"`
	Limit int `qstring:"limit"`
}

// Query is the http request query struct.
type Query struct {
	Names    []string
	IDs      []int
	PageInfo *PagingParams
	Created  time.Time
	Modified time.Time
}
```

## Additional Notes

- All Timestamps are assumed to be in RFC3339 format
- A struct field tag of `qstring` is supported and supports all of the features you've come to know and love from Go
  (un)marshalers.
  - A field tag with a value of `qstring:"-"` instructs `qstring` to ignore the field.
  - A field tag with an the `omitempty` option set will be ignored if the field being marshaled has a zero value.
    `qstring:"name,omitempty"`
  - A slice field tag with the `comma` option set will produce a comma-separated list of items in a single query
    parameter instead of multiple instances of the parameter. `qstring:"names,comma"` will produce `?names=a,b,c`
    instead of `?names=a&names=b&names=c`.

### Custom Fields

In order to facilitate more complex queries `qstring` also provides some custom fields to save you a bit of headache
with custom marshal/unmarshaling logic. Currently the following custom fields are provided:

- `qstring.ComparativeTime` - Supports timestamp query parameters with optional logical operators (<, >, <=, >=) such as
  `?created<=2006-01-02T15:04:05Z`

## Benchmarks

```
BenchmarkUnmarshall-4 	  500000	      2711 ns/op	     448 B/op	      23 allocs/op
BenchmarkRawPLiteral-4	 1000000	      1675 ns/op	     448 B/op	      23 allocs/op
ok  	github.com/Southclaws/qstring	3.163s
```
