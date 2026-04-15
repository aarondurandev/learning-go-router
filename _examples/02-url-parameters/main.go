package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Can't find it")
	})
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "home")
	})

	m.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello there")
	})
	m.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, %s", router.URLParam(r, "name"))
	})
	m.Get("/hello/{name}/{age}", func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")
		age := router.URLParam(r, "age")
		fmt.Fprintf(w, "Hello, %s. You're %s years old", name, age)
	})
	http.ListenAndServe(":8080", m)
}
