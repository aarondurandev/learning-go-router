// This example demonstrates wildcard routes. A trailing * segment matches
// the rest of the path. The captured value is retrieved with URLParam(r, "*").
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		file := router.URLParam(r, "*")
		if file != "" {
			fmt.Fprintf(w, "Requested file: %s", file)
		}
	})
	http.ListenAndServe(":8080", m)
}
