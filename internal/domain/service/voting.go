package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/yvv4git/task-voting/internal/domain/repository"
	"github.com/yvv4git/task-voting/internal/interfaces/web"
)

type VotingRepository interface {
	List(ctx context.Context, r *repository.ListVotingRequest) (*repository.ListVotingResponse, error)
	CreateVoting(context.Context, *repository.CreateVotingParams) (*repository.CreateVotingResult, error)
	UpdateVoting(context.Context, *repository.UpdateVotingParams) error
	DeleteVoting(context.Context, *repository.DeleteVotingParams) error
	MakeChoice(context.Context, *repository.MakeChoiceParams) error
}

type SubscriptionProcessor interface {
	Broadcast(message []byte)
}

type Voting struct {
	logger       *slog.Logger
	repo         VotingRepository
	subscription SubscriptionProcessor
}

func NewVoting(logger *slog.Logger, repo VotingRepository, subscription SubscriptionProcessor) *Voting {
	return &Voting{
		logger:       logger,
		repo:         repo,
		subscription: subscription,
	}
}

func (v *Voting) List(ctx context.Context, r *web.ListVotingRequest) (*web.ListVotingResponse, error) {
	data, err := v.repo.List(ctx, &repository.ListVotingRequest{
		Limit:  r.Limit,
		Offset: r.Offset,
	})
	if err != nil {
		return nil, err
	}

	var result = make([]web.VotingItem, 0, len(data.Items))
	for _, item := range data.Items {
		var dstItem web.VotingItem
		dstItem.ID = item.ID
		dstItem.Name = item.Name
		dstItem.Description = item.Description
		dstItem.CreatedAt = item.CreatedAt
		dstItem.StartAt = item.StartAt
		dstItem.EndAt = item.EndAt
		if len(item.Invariance) > 0 {
			var dstScoreItems = make([]web.InvarianceScore, 0, len(item.Invariance))
			for _, invariance := range item.Invariance {
				var dstInvariance web.InvarianceScore
				dstInvariance.ID = invariance.ID
				dstInvariance.Name = invariance.Name
				dstInvariance.Score = invariance.Score
				dstScoreItems = append(dstScoreItems, dstInvariance)
			}
			dstItem.Invariance = dstScoreItems
		}

		result = append(result, dstItem)
	}

	return &web.ListVotingResponse{
		Items: result,
	}, nil
}

func (v *Voting) CreateVoting(ctx context.Context, r *web.CreateVotingRequest) (*web.CreateVotingResponse, error) {
	result, err := v.repo.CreateVoting(ctx, &repository.CreateVotingParams{
		Name:        r.Name,
		Description: r.Description,
		StartAt:     r.StartAt,
		EndAt:       r.EndAt,
		Invariance:  r.Invariance,
	})
	if err != nil {
		return nil, err
	}

	return &web.CreateVotingResponse{
		ID: result.ID,
	}, nil
}

func (v *Voting) UpdateVoting(ctx context.Context, r *web.UpdateVotingRequest) error {
	return v.repo.UpdateVoting(ctx, &repository.UpdateVotingParams{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		StartAt:     r.StartAt,
		EndAt:       r.EndAt,
		Invariance:  r.Invariance,
	})
}

func (v *Voting) DeleteVoting(ctx context.Context, r *web.DeleteVotingRequest) error {
	return v.repo.DeleteVoting(ctx, &repository.DeleteVotingParams{
		ID: r.ID,
	})
}

func (v *Voting) MakeChoice(ctx context.Context, r *web.MakeChoiceRequest) error {
	const (
		limit  = 100
		offset = 0
	)
	defer func() {
		listVoting, err := v.repo.List(ctx, &repository.ListVotingRequest{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			// TODO: log
			return
		}

		payload, err := json.Marshal(listVoting)
		if err != nil {
			// TODO: log
			return
		}

		v.subscription.Broadcast(payload)
	}()

	return v.repo.MakeChoice(ctx, &repository.MakeChoiceParams{
		InvarianceID: r.InvarianceID,
		UserID:       r.UserID,
	})
}
