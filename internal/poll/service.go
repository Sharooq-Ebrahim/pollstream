package poll

type PollService struct {
	repo *PollRepository
	hub  *Hub
}

func NewPollService(repo *PollRepository, hub *Hub) *PollService {
	return &PollService{repo: repo, hub: hub}
}

func (ps *PollService) CreatePoll(poll *Poll) error {

	err := ps.repo.CreatePoll(poll)
	if err != nil {
		return err
	}

	ps.hub.Broadcast(poll)
	return nil
}

func (ps *PollService) GetPollByID(id string) (*Poll, error) {
	return ps.repo.GetPollByID(id)
}
