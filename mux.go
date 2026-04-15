package router

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is an unexported type for context keys in this package.
// Using a named type prevents collisions with keys from other packages.
type contextKey string

// paramsKey is the context key under which URL parameters are stored.
const paramsKey contextKey = "params"

// Route represents a single registered route, holding the HTTP method,
// URL pattern, and the handler to invoke on a match.
type Route struct {
	method  string
	pattern string
	handler http.Handler
}

// Mux is the router. It stores registered routes and dispatches
// incoming requests to the appropriate handler.
type Mux struct {
	routes          []Route
	notFoundHandler http.HandlerFunc
}

// compile-time check that *Mux implements Router.
var _ Router = (*Mux)(nil)

// Method registers a handler for the given HTTP method and pattern.
// The method is normalized to uppercase before storing.
// All other registration methods delegate to this one.
func (mx *Mux) Method(method, pattern string, handler http.Handler) {
	parsedMethod := strings.ToUpper(method)
	newRoute := Route{
		method:  parsedMethod,
		pattern: pattern,
		handler: handler,
	}
	mx.routes = append(mx.routes, newRoute)
}

// MethodFunc registers a handler function for the given HTTP method and pattern.
func (mx *Mux) MethodFunc(method, pattern string, handlerFn http.HandlerFunc) {
	mx.Method(method, pattern, handlerFn)
}

// Handle registers a handler for the given pattern, matching any HTTP method.
// Internally the method is stored as an empty string to signal "any method".
func (mx *Mux) Handle(pattern string, handler http.Handler) {
	mx.Method("", pattern, handler)
}

// HandleFunc registers a handler function for the given pattern, matching any HTTP method.
func (mx *Mux) HandleFunc(pattern string, handlerFn http.HandlerFunc) {
	mx.Method("", pattern, handlerFn)
}

func (mx *Mux) Get(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("GET", pattern, handlerFn)
}
func (mx *Mux) Delete(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("DELETE", pattern, handlerFn)
}
func (mx *Mux) Post(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("POST", pattern, handlerFn)
}
func (mx *Mux) Patch(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("PATCH", pattern, handlerFn)
}
func (mx *Mux) Put(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("PUT", pattern, handlerFn)
}

// NotFound sets a custom handler for requests that match no route.
// If not set, the default net/http not-found handler is used.
func (mx *Mux) NotFound(handlerFn http.HandlerFunc) {
	mx.notFoundHandler = handlerFn
}

// ServeHTTP dispatches the request to the matching route's handler.
// Returns 405 if the path matches but the method does not.
// Returns 404 if no route matches at all.
func (mx *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var matchFound bool = false
	for _, route := range mx.routes {
		matched, params := matchPath(route.pattern, r.URL.Path)
		if matched && (route.method == r.Method || route.method == "") {
			ctx := context.WithValue(r.Context(), paramsKey, params)
			route.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		} else if matched && (route.method != r.Method) {
			matchFound = true
		}

	}
	if matchFound == true {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		if mx.notFoundHandler != nil {
			mx.notFoundHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}
// URLParam returns the value of the URL parameter with the given key
// from the request context. Returns an empty string if not found.
func URLParam(r *http.Request, key string) string {
	params, ok := r.Context().Value(paramsKey).(map[string]string)
	if !ok {
		return ""
	}
	return params[key]
}

// NewMux creates and returns a new Mux instance.
func NewMux() *Mux {
	return &Mux{}
}

// matchPath compares a route pattern against a request path.
// Segments wrapped in {} are treated as parameters and captured into the returned map.
// Returns false and nil if the path does not match the pattern.
func matchPath(pattern string, path string) (bool, map[string]string) {
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(path, "/")
	if len(patternSegments) != len(pathSegments) {
		return false, nil
	}
	params := make(map[string]string)
	for i, seg := range patternSegments {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			name := seg[1 : len(seg)-1]
			params[name] = pathSegments[i]
		} else if seg != pathSegments[i] {
			return false, nil
		}
	}
	return true, params
}
