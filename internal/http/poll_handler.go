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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(poll)
	w.WriteHeader(http.StatusOK)

}

func (ph *PollHandler) GetAllPolls(w http.ResponseWriter, r *http.Request) {

	polls, err := ph.service.GetAllPolls()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(polls)
	w.WriteHeader(http.StatusOK)

}

func (ph *PollHandler) Vote(w http.ResponseWriter, r *http.Request) {

	type vote struct {
		PollID   string `json:"poll_id"`
		OptionID string `json:"option_id"`
	}

	var v vote
	err := json.NewDecoder(r.Body).Decode(&v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ph.service.Vote(v.PollID, v.OptionID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
