package poll

type PollService struct {
	repo *PollRepository
}

func NewPollService(repo *PollRepository) *PollService {
	return &PollService{repo: repo}
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
