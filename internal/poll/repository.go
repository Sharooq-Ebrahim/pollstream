package poll

import (
	"database/sql"
)

type PollRepository struct {
	db *sql.DB
}

func NewPollRepository(db *sql.DB) *PollRepository {
	return &PollRepository{db: db}
}

func (r *PollRepository) CreatePoll(poll *Poll) error {

	txn, err := r.db.Begin()

	if err != nil {
		return err
	}

	_, err = txn.Exec("INSERT INTO polls (id,question) VALUES ($1,$2) RETURNING id", poll.ID, poll.Question)

	if err != nil {
		txn.Rollback()
		return err
	}

	for _, option := range poll.Options {

		_, err := txn.Exec("INSERT INTO options (id,poll_id,text,votes) VALUES ($1,$2,$3,$4)", option.ID, poll.ID, option.Text, option.Votes)

		if err != nil {
			txn.Rollback()
			return err
		}

	}

	err = txn.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (r *PollRepository) GetPollByID(id string) (*Poll, error) {

	var p Poll

	err := r.db.QueryRow("SELECT id,question FROM polls WHERE id=$1", id).Scan(&p.ID, &p.Question)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query("SELECT id,text,votes from options where poll_id=$1", p.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var option Option

		err := rows.Scan(&option.ID, &option.Text, &option.Votes)

		if err != nil {
			return nil, err
		}

		p.Options = append(p.Options, option)
	}

	return &p, nil
}
