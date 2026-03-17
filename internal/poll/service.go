package poll

type PollService struct {
	repo *PollRepository
	hub  *Hub
}

func NewPollService(repo *PollRepository, hub *Hub) *PollService {
	return &PollService{repo: repo, hub: hub}
}

func (ps *PollService) CreatePoll(poll *Poll) error {

	return ps.repo.CreatePoll(poll)

}

func (ps *PollService) GetPollByID(id string) (*Poll, error) {
	return ps.repo.GetPollByID(id)
}

func (ps *PollService) GetAllPolls() ([]Poll, error) {
	return ps.repo.GetAllPolls()
}

func (ps *PollService) Vote(pollID string, optionID string) error {
	err := ps.repo.Vote(pollID, optionID)
	if err != nil {
		return err
	}

	poll, err := ps.repo.GetPollByID(pollID)
	if err == nil {
		ps.hub.Broadcast(poll)
	}

	return nil
}
