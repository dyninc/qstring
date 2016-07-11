# qstring
The package provides an easy way to marshal and unmarshal query string data to
and from structs.

## Installation
```bash
$ go get github.com/dyninc/qstring
```

## Examples

### Basic Unmarshal

```go
package main

import (
	"net/http"

	"github.com/dyninc/qstring"
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

The above example will unmarshal the query string from an http.Request and
unmarshal it into the provided struct. This means that a query of
`` would be unmarshaled into a struct similar
to the following:

```go
Query{
	Names: []string{"foo", "bar"},
	Limit: 50,
	Page: 1,
}
```

### Marshalling
`qstring` exposes two methods of Marshaling structs into Query parameters, one
will Marshal the provided struct into a raw query string, the other will Marshal
a struct into a `url.Values` type. Some Examples of both follow.

### Marshal Raw Query String
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/dyninc/qstring"
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
  // Output: names=foo&names=bar&limit=50&page=1
}
```

### Marshal url.Values
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/dyninc/qstring"
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
  q, err := qstring.MarshalValues(query)
	fmt.Println(q)
  // Output: map[names:[foo, bar] limit:[50] page:[1]]
}
```

### Nested
In the same spirit as other Unmarshaling libraries, `qstring` allows you to
Marshal/Unmarshal nested structs

```go
package main

import (
	"net/http"

	"github.com/dyninc/qstring"
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

## Additional Notes
* All Timestamps are assumed to be in RFC3339 format
* A struct field tag of `qstring` is supported and. This supports all of the
functionality you've come to know and love. Ie a value of `"-"` instructs `qstring`
to ignore the field.

### Complex Structures
Again, in the spirit of other Unmarshaling libraries, `qstring` allows for some
more complex types, such as pointers and time.Time fields. A more complete
example might look something like the following code snippet

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
