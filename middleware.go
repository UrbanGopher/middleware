package middleware

import (
	"net/http"
)

// Middleware - A function type for defining middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Decorator - Interface for defining a middleware wrapping mechanism. This
// interface exists for the purpose of testing, but, of course, can be used
// at your discretion.
type Decorator interface {
	// Append - Add Middleware to the beginning of the decorator chain. This
	// Middleware will wrap all existing middleware.
	Append(Middleware) Decorator

	// Prepend - Add Middleware to the end of the decorator chain. This Middleware
	// will be wrapped by all existing Middleware.
	Prepend(Middleware) Decorator

	// WrapHandler - Wraps the http.HandlerFunc with the Middleware added via the
	// Append and Prepend methods
	WrapHandler(http.HandlerFunc) http.HandlerFunc
}

// Wrapper - Concrete implementation of the Decorator interface
type Wrapper struct {
	middlewares []Middleware
}

// Append - Concrete implementation of the Decorator.Append method
func (w *Wrapper) Append(mw Middleware) Decorator {
	w.middlewares = append([]Middleware{mw}, w.middlewares...)
	return w
}

// Prepend - Concrete implementation of the Decorator.Prepend method
func (w *Wrapper) Prepend(m Middleware) Decorator {
	w.middlewares = append(w.middlewares, m)
	return w
}

// WrapHandler - Concrete implementation of the Decorator.WrapHandler method
func (w *Wrapper) WrapHandler(h http.HandlerFunc) http.HandlerFunc {
	var mw Middleware
	for _, mw = range w.middlewares {
		h = mw(h)
	}
	return h
}
