package handlers

import (
	"encoding/json"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/rs/zerolog/log"
	"net/http"
)

// List In case some response type or sth similar is needed in the future
type List struct {
	res         []*entity.Task
	TaskService interfaces.ITaskService
}

// @Summary list tasks
// @Description  list the existing tasks
// @Produce json
// @Success 201 {array} entity.Task
// @Failure 405,400,500
// @Router /tasks [get]
//
// ServeHTTP implements the handler interface to handle creating the tasks
func (l List) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := l.TaskService.Get()
	l.res = tasks
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(l.res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to write response")
		return
	}
	return
}
