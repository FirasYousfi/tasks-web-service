package handlers

import (
	"encoding/json"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// Create is the struct used for that will implement the Handler interface for the creation
type Create struct {
	req         entity.TaskDescription
	res         entity.Task
	TaskService interfaces.ITaskService
}

// @Summary create a task
// @Description  add a new task to the tasks list
// @Produce json
// @Accept	json
// @Param   task  body  entity.TaskDescription  true  "New task"
// @Success 201 {object} entity.Task
// @Failure 405,400,500
// @Router /tasks [post]
func (c Create) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body) // we use decoder instead of unmarshall
	err := decoder.Decode(&c.req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed close body")
		}
	}(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //bed request because body is not json
		log.Error().Err(err).Msg("failed to decode body")
		return
	}

	response, err := c.TaskService.Create(&c.req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to create task")
		return
	}

	c.res = *response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(c.res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to write response")
		return
	}
	return
}

///***INFOS****////
/*
Why defer?
In the Go sample, if the "ParseResponse" call returns an error (and the function returns "early"), the "defer"
will ensure that the response body will be closed. If the "ParseResponse" call is successful, then the "defer" will
ensure that the response body will be closed.
In both scenarios, the response body will be closed whether the operation was successful or whether it failed.
*/
