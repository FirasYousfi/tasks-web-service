package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

// @Summary get a task
// @Description  get a specific task by its ID
// @Produce json
// @Param id path string true "task ID"
// @Success 200 {object} entity.Task
// @Failure 405,400,500
// @Router /tasks/{id} [get]
//
// ServeHTTP implements the handler interface to handle getting a task by ID
func (g GetTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["taskId"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg("task ID not provided in path")
		return
	}
	task, err := g.Service.GetTaskByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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

func (g GetCollection) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["collectionId"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg("collection ID not provided in path")
		return
	}
	collection, err := g.Service.GetCollectionsByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msgf("failed to find task with id %s", id)
		return
	}
	g.res = *collection
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
