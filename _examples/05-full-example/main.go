// This example combines all router features: global middleware, route groups
// with their own middleware stacks, URL parameters, and a custom 404 handler.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.Use(globalMiddleware)
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404 - not found")
	})
	m.Group("/api/v1", func(r router.Router) {
		r.Use(firstGroupMiddleware)
		r.Get("/users", usersHandlerV1)
		r.Get("/users/{id}", getUserHandlerV1)
		r.Post("/users", createUserHandlerV1)
	})
	m.Group("/api/v2", func(r router.Router) {
		r.Use(secondGroupMiddleware)
		r.Get("/users", usersHandlerV2)
		r.Get("/users/{id}", getUserHandlerV2)
		r.Post("/users", createUserHandlerV2)
	})
	http.ListenAndServe(":8080", m)
}

func usersHandlerV1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "users endpoint v1")
}

func getUserHandlerV1(w http.ResponseWriter, r *http.Request) {
	id := router.URLParam(r, "id")
	fmt.Fprintf(w, "v1 user: %s", id)
}

func createUserHandlerV1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "v1 create user")
}

func usersHandlerV2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "users endpoint v2")
}

func getUserHandlerV2(w http.ResponseWriter, r *http.Request) {
	id := router.URLParam(r, "id")
	fmt.Fprintf(w, "v2 user: %s", id)
}

func createUserHandlerV2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "v2 create user")
}

func globalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func firstGroupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("v1 middleware")
		next.ServeHTTP(w, r)
	})
}

func secondGroupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("v2 middleware")
		next.ServeHTTP(w, r)
	})
}
