# qstring
The package provides an easy way to marshal and unmarshal query string data to
and from structs.

# Installation
```bash
$ go get github.com/dyninc/qstring
```

# Example

## Literal

```go
package main

import (
	"net/http"

	"github.com/dyninc/qstring"
)

// Query is the request query struct.
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

## Nested
In the same spirit as other Unmarshalling libraries, `qstring` allows you to
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

// Req is the request query struct.
type Query struct {
  Names    []string
  PageInfo PagingParams
}
```

## Complex Structures
Again, in the spirit of other Unmarshalling libraries, `qstring` allows for some
more complex types, such as pointers and time.Time fields. *Note: All Timestamps
are assumed to be in RFC3339 format*.

```go
package main

import (
	"net/http"
  "time"

	"github.com/dyninc/qstring"
)

// PagingParams represents common pagination information for query strings
type PagingParams struct {
  Page int
  Limit int
}

// Req is the request query struct.
type Query struct {
  Names    []string
  PageInfo *PagingParams
  Created  time.TIme
  Modified time.Time
}
```
