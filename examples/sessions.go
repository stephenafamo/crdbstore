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
