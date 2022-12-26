package router

import (
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/adapters/web/handlers"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/k8s"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
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

	return r
}
