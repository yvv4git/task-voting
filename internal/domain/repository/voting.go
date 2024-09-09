package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yvv4git/task-voting/internal/infrastructure"

	sq "github.com/Masterminds/squirrel"
)

const (
	tbVoting           = "voting"
	tbVotingInvariance = "voting_invariance"
)

type Voting struct {
	db *pgxpool.Pool
}

func NewVoting(db *pgxpool.Pool) *Voting {
	return &Voting{
		db: db,
	}
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

type ResultID struct {
	ID uuid.UUID `db:"id"`
}

func (v *Voting) CreateVoting(ctx context.Context, p *CreateVotingParams) (*CreateVotingResult, error) {
	tx, err := v.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	resultID, err := v.createVoting(ctx, tx, p)
	if err != nil {
		return nil, fmt.Errorf("create voting: %w", err)
	}

	if err = v.addInvariance(ctx, tx, addInvarianceParams{
		votingID:   resultID.ID,
		invariance: p.Invariance,
	}); err != nil {
		return nil, fmt.Errorf("add invariance: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &CreateVotingResult{
		ID: resultID.ID,
	}, nil
}

func (v *Voting) createVoting(ctx context.Context, tx pgx.Tx, p *CreateVotingParams) (*ResultID, error) {
	b := sq.Insert(tbVoting).
		Columns("name", "description", "started_at", "ended_at").
		Values(p.Name, p.Description, p.StartAt, p.EndAt).
		Suffix("RETURNING id")

	stmt, args, err := b.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build statement: %w", err)
	}

	resultIDPtr, err := infrastructure.FetchRow[ResultID](ctx, tx, stmt, args...)
	if err != nil {
		return nil, err
	}

	return resultIDPtr, nil
}

type addInvarianceParams struct {
	votingID   uuid.UUID
	invariance []string
}

func (v *Voting) addInvariance(ctx context.Context, tx pgx.Tx, p addInvarianceParams) error {
	b := sq.Insert(tbVotingInvariance).
		Columns("voting_id", "name")

	for _, item := range p.invariance {
		b = b.Values(p.votingID, item)
	}

	stmt, args, err := b.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	if err = infrastructure.Execute(ctx, tx, stmt, args...); err != nil {
		return err
	}

	return nil
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
