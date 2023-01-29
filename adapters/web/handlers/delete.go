package handlers

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

// @Summary delete a collection
// @Description  delete a collection from the list
// @Param collectionId path string true "collection ID"
// @Success 200
// @Failure 405,400,500
// @Router /tasks/{id} [delete]
//
// ServeHTTP implements the handler interface to handle deleting the collections
func (d DeleteCollection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := mux.Vars(r)["collectionId"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg("task ID not provided in request path")
		return
	}
	_, err := d.Service.GetCollectionsByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Error().Msgf("error: %s, occurred when getting task with ID %s", err.Error(), id)
		return
	}

	err = d.Service.DeleteCollectionByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to delete instance")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
