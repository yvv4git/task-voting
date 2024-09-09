package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Voting struct {
	// db
}

func NewVoting() *Voting {
	return &Voting{}
}

type (
	ListVotingRequest struct {
		Limit  int64
		Offset int64
	}

	InvarianceScore struct {
		Name  string
		Score int64
	}

	VotingItem struct {
		ID          uuid.UUID
		Name        string
		Description string
		CreatedAt   time.Time
		StartAt     time.Time
		EndAt       time.Time
		Invariance  []InvarianceScore
	}

	ListVotingResponse struct {
		Items []VotingItem
	}
)

func (v *Voting) List(ctx context.Context, r *ListVotingRequest) (*ListVotingResponse, error) {
	// Example
	response := &ListVotingResponse{
		Items: []VotingItem{
			{
				ID:          uuid.New(),
				Name:        "Voting 1",
				Description: "This is the first voting item",
				CreatedAt:   time.Now(),
				StartAt:     time.Now().Add(time.Hour * 24),
				EndAt:       time.Now().Add(time.Hour * 48),
				Invariance: []InvarianceScore{
					{
						Name:  "Invariance 1",
						Score: 90,
					},
					{
						Name:  "Invariance 2",
						Score: 85,
					},
				},
			},
			{
				ID:          uuid.New(),
				Name:        "Voting 2",
				Description: "This is the second voting item",
				CreatedAt:   time.Now(),
				StartAt:     time.Now().Add(time.Hour * 48),
				EndAt:       time.Now().Add(time.Hour * 72),
				Invariance: []InvarianceScore{
					{
						Name:  "Invariance 1",
						Score: 75,
					},
					{
						Name:  "Invariance 2",
						Score: 80,
					},
				},
			},
		},
	}

	return response, nil
}

type (
	CreateVotingParams struct {
		Name        string
		Description string
		StartAt     time.Time
		EndAt       time.Time
		Invariance  []string
	}

	CreateVotingResult struct {
		ID uuid.UUID
	}
)

func (v *Voting) CreateVoting(_ context.Context, p *CreateVotingParams) (*CreateVotingResult, error) {
	// TODO: implement
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &CreateVotingResult{
		ID: id,
	}, nil
}

type UpdateVotingParams struct {
	Name        string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	Invariance  []string
}

func (v *Voting) UpdateVoting(_ context.Context, p *UpdateVotingParams) error {
	// TODO: implement
	return nil
}

type DeleteVotingParams struct {
	ID uuid.UUID
}

func (v *Voting) DeleteVoting(_ context.Context, p *DeleteVotingParams) error {
	// TODO: implement
	return nil
}

type MakeChoiceParams struct {
	InvarianceID uuid.UUID `json:"invarianceId"`
}

func (v *Voting) MakeChoice(_ context.Context, p *MakeChoiceParams) error {
	// TODO: implement
	return nil
}
