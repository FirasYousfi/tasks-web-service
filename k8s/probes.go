package k8s

import (
	"gorm.io/gorm"
	"net/http"
)

type Liveness struct {
}

type Readiness struct {
	DB *gorm.DB
}

// ServeHTTP defines the handling of liveness probe, checks just if app is alive
func (l Liveness) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ServeHTTP defines the handling of readiness probe, checks if app is ready for requests by seeing if DB is set and working.
func (r Readiness) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// if db is completely nil we cannot be ready
	if r.DB == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	db, err := r.DB.DB()
	if err != nil || db.Ping() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
