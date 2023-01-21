package router

import "net/http"

type Middleware func(http.Handler) http.Handler

// attachMiddleware chains multiple middleware function to the provided http handler.
// it is preferred over r.Use(basicAuth) of the gorilla/mux lib because it easily allows different middlewares for different routes
func attachMiddleware(h http.Handler, m ...Middleware) http.Handler {

	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped

}
