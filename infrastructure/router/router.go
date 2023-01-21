package router

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/adapters/web/handlers"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/config"
	"github.com/FirasYousfi/tasks-web-servcie/k8s"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

const basePath = "/v1/api"

func SetupRoutes(service interfaces.ITaskService, db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	r.Handle(fmt.Sprintf("%s/tasks", basePath), &handlers.Create{TaskService: service}).Methods("POST")
	r.Handle(fmt.Sprintf("%s/tasks", basePath), &handlers.List{TaskService: service}).Methods("GET")
	r.Handle(fmt.Sprintf("%s/tasks/{id}", basePath), &handlers.Delete{TaskService: service}).Methods("DELETE")
	r.Handle(fmt.Sprintf("%s/tasks/{id}", basePath), &handlers.Get{TaskService: service}).Methods("GET")
	r.Handle(fmt.Sprintf("%s/tasks/{id}", basePath), &handlers.Update{TaskService: service}).Methods("PATCH")
	r.Handle(fmt.Sprintf("%s/tasks/{id}", basePath), &handlers.Update{TaskService: service}).Methods("PUT")

	// liveness and readiness probes
	r.Handle(fmt.Sprintf("/healthz"), &k8s.Liveness{}).Methods("GET")
	r.Handle(fmt.Sprintf("/readyz"), &k8s.Readiness{DB: db}).Methods("GET")
	// using middleware with gorilla mux
	r.Use(basicAuth)
	return r
}

// To use middleware with the r.Use(MiddlewareFunc) provided by gorilla mux, the signature needs to be: type MiddlewareFunc func(http.Handler) http.Handler
func basicAuth(next http.Handler) http.Handler {
	//if 'f' is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request header. If the Authentication header is not present or is invalid we know through the 'ok' variable
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for the provided and expected usernames and passwords so that we ensure they have the same length
			// that way the comparison with subtle.ConstantTimeCompare would always take the same time
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(config.Config.Auth.Username))
			expectedPasswordHash := sha256.Sum256([]byte(config.Config.Auth.Password))

			// Use the subtle.ConstantTimeCompare() function to check if the provided username and password hashes equal the expected ones.
			// Importantly, we should to do the work to evaluate both the username and password before checking the return values to avoid leaking information.
			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			// If the username and password are correct, then call
			// the next handler in the chain. Make sure to return
			// afterwards, so that none of the code below is run.
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
