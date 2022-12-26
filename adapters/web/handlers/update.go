package handlers

import (
	"encoding/json"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// Update represents the struct that implement the handler for update
type Update struct {
	req         entity.TaskDescription
	res         entity.Task
	TaskService interfaces.ITaskService
}

// @Summary update a task
// @Description  update a task by ID
// @Param id path string true "task ID"
// @Param   task  body  entity.TaskDescription  true  "New task description"
// @Produce json
// @Accept	json
// @Success 200 {object} entity.Task
// @Failure 405,400,500
// @Router /tasks/{id} [put]
// @Router /tasks/{id} [patch]
//
// ServeHTTP implements the handler interface to handle updating the tasks
func (u Update) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Both requests have the same path, so they should have the same handler depending on the request method
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u.req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed close body")
		}
	}(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(err).Msg("failed to decode body")
		return
	}
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg("task ID not provided in path")
		return
	}
	var response *entity.Task
	if r.Method == http.MethodPut { //PUT here
		response, err = u.TaskService.UpdateFully(&u.req, id)
	} else { //PATCH here
		response, err = u.TaskService.UpdatePartial(&u.req, id)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to update task")
		return
	}
	u.res = *response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(u.res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to write response")
		return
	}
	return
}
