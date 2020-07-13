# crdbstore

A session store backend for [gorilla/sessions](http://www.gorillatoolkit.org/pkg/sessions) - [src](https://github.com/gorilla/sessions).

## Installation

    make get-deps

## Documentation

Available on [godoc.org](http://pkg.go.dev/github.com/stephenafamo/cockroachstore).

See http://www.gorillatoolkit.org/pkg/sessions for full documentation on underlying interface.

### Example

[embedmd]:# (examples/sessions.go)
```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/stephenafamo/crdbstore"
)

var store *crdbstore.CrDBStore

func main() {
	var err error
	// Fetch new store.
	store, err = crdbstore.NewCrDBStore("postgres://user:password@127.0.0.1:5432/database?sslmode=verify-full", []byte("secret-key"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	http.HandleFunc("/", ExampleHandler)
	http.ListenAndServe(":8080", nil)
}

// ExampleHandler is an example that displays the usage of PGStore.
func ExampleHandler(w http.ResponseWriter, r *http.Request) {

	// Get a session.
	session, err := store.Get(r, "session-key")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Add a value.
	session.Values["foo"] = "bar"

	// Save.
	if err = session.Save(r, w); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	// Delete session.
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}
}
```

## Thanks

This driver is mostly copied from the postgres driver which no longer works for cockroachDB.

* [pgstore](https://github.com/antonlindstrom/pgstore)

What makes this backend different is that it's for CockroachDB. 
