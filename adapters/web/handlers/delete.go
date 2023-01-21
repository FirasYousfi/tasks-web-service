package handlers

import (
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Delete struct {
	TaskService interfaces.ITaskService
}

// @Summary delete a task
// @Description  delete a task from the list
// @Param id path string true "task ID"
// @Success 200
// @Failure 405,400,500
// @Router /tasks/{id} [delete]
//
// ServeHTTP implements the handler interface to handle deleting the tasks
func (d Delete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg("task ID not provided in request path")
		return
	}
	_, err := d.TaskService.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Error().Msgf("error: %s, occurred when getting task with ID %s", err.Error(), id)
		return
	}

	err = d.TaskService.DeleteByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to delete instance")
		return
	}
	// why we are using http.StatusNoContent: https://stackoverflow.com/questions/2342579/http-status-code-for-update-and-delete#:~:text=For%20a%20DELETE%20request%3A%20HTTP,but%20not%20fully%20applied%20yet.
	w.WriteHeader(http.StatusNoContent)
	return
}
