package handlers

import (
	"encoding/json"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Get In case some response type or sth similar is needed in the future
type Get struct {
	res         entity.Task
	TaskService interfaces.ITaskService
}

// @Summary get a task
// @Description  get a specific task by its ID
// @Produce json
// @Param id path string true "task ID"
// @Success 200 {object} entity.Task
// @Failure 405,400,500
// @Router /tasks/{id} [get]
//
// ServeHTTP implements the handler interface to handle getting a task by ID
func (g Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg("task ID not provided in path")
		return
	}
	task, err := g.TaskService.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msgf("failed to find task with id %s", id)
		return
	}
	g.res = *task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(g.res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to write response")
		return
	}
	return
}
