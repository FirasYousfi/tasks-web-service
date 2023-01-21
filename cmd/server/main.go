package main

import (
	"github.com/FirasYousfi/tasks-web-servcie/adapters/persistence/repository"
	"github.com/FirasYousfi/tasks-web-servcie/application/service"
	"github.com/FirasYousfi/tasks-web-servcie/config"
	"github.com/FirasYousfi/tasks-web-servcie/infrastructure/database"
	"github.com/FirasYousfi/tasks-web-servcie/infrastructure/router"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// @title           Tasks Service API
// @version         1.0
// @description     This is the documentation for the tasks-service-api.

// @contact.name   Firas Yousfi
// @contact.email  firas.yousfi144@gmail.com

// @schemes http
// @host      localhost:8080
// @BasePath  /v1/api
func main() {
	config.BuildConfig()
	err := database.InitializeDB()
	if err != nil {
		log.Fatalf("error Initializing database: %v", err)
	}
	r := SetupHandlers(database.DB.GetDBConn())
	log.Println("Serving on port 8080:")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// SetupHandlers here is where all the dependency injection stuff happens.
func SetupHandlers(db *gorm.DB) *mux.Router {
	repo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(repo)
	r := router.SetupRoutes(taskService, db)
	return r
}
