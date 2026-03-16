package api

import (
	"encoding/json"
	"net/http"
	"pollstream/internal/poll"
)

type PollHandler struct {
	service *poll.PollService
}

func NewPollHandler(service *poll.PollService) *PollHandler {
	return &PollHandler{service: service}

}

func (ph *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {

	var poll poll.Poll

	err := json.NewDecoder(r.Body).Decode(&poll)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ph.service.CreatePoll(&poll)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (ph *PollHandler) GetPollByID(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	poll, err := ph.service.GetPollByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(poll)

	w.WriteHeader(http.StatusOK)

}
