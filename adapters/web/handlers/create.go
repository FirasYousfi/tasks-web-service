package handlers

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// @Summary create a collection
// @Description  add a new collection to the collections list
// @Produce json
// @Accept	json
// @Param   collection  body  entity.CollectionDescription  true  "New collection"
// @Success 201 {object} entity.Collection
// @Failure 405,400,500
// @Router /collections [post]
//
// ServeHTTP implements the handler interface to handle deleting the collections
func (c CreateCollection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	response, err := c.Service.CreateCollection(&c.req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to create collection")
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
