package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VotingService interface {
	List(ctx context.Context, r *ListVotingRequest) (*ListVotingResponse, error)
	CreateVoting(ctx context.Context, r *CreateVotingRequest) (*CreateVotingResponse, error)
	UpdateVoting(ctx context.Context, r *UpdateVotingRequest) error
	DeleteVoting(ctx context.Context, r *DeleteVotingRequest) error
	MakeChoice(ctx context.Context, r *MakeChoiceRequest) error
}

type VotingHandler struct {
	votingService VotingService
}

func NewVotingHandler(votingService VotingService) *VotingHandler {
	return &VotingHandler{
		votingService: votingService,
	}
}

func (v *VotingHandler) RegisterHandlers(router *gin.Engine) {
	router.GET("/voting", v.ListVoting)
	router.POST("/voting", v.CreateVoting)
	router.PUT("/voting/:id", v.UpdateVoting)
	router.DELETE("/voting/:id", v.DeleteVoting)
	router.PUT("/voting/:id/choice", v.MakeChoice)
}

type (
	ListVotingRequest struct {
		Limit  int64 `form:"limit"`
		Offset int64 `form:"offset"`
	}

	InvarianceScore struct {
		Name  string
		Score int64
	}

	VotingItem struct {
		ID          uuid.UUID         `json:"id"`
		Name        string            `json:"name"`
		Description string            `json:"description"`
		CreatedAt   time.Time         `json:"created_at"`
		StartAt     time.Time         `json:"startAt"`
		EndAt       time.Time         `json:"endAt"`
		Invariance  []InvarianceScore `json:"invariance"`
	}

	ListVotingResponse struct {
		Items []VotingItem `json:"items"`
	}
)

func (v *VotingHandler) ListVoting(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	result, err := v.votingService.List(c.Request.Context(), &ListVotingRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"votings": result})
}

type (
	CreateVotingRequest struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		StartAt     time.Time `json:"startAt"`
		EndAt       time.Time `json:"endAt"`
		Invariance  []string  `json:"invariance"`
	}

	CreateVotingResponse struct {
		ID uuid.UUID `json:"id"`
	}
)

func (v *VotingHandler) CreateVoting(c *gin.Context) {
	var request CreateVotingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := v.votingService.CreateVoting(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	response := CreateVotingResponse{
		ID: result.ID,
	}

	c.JSON(http.StatusCreated, response)
}

type UpdateVotingRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"startAt"`
	EndAt       time.Time `json:"endAt"`
	Invariance  []string  `json:"invariance"`
}

func (v *VotingHandler) UpdateVoting(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	fmt.Println("id: ", id)

	var request UpdateVotingRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = v.votingService.UpdateVoting(c.Request.Context(), &UpdateVotingRequest{
		Name:        request.Name,
		Description: request.Description,
		StartAt:     request.StartAt,
		EndAt:       request.EndAt,
		Invariance:  request.Invariance,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "voting updated"})
}

type DeleteVotingRequest struct {
	ID uuid.UUID `json:"id"`
}

func (v *VotingHandler) DeleteVoting(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	fmt.Println("id: ", id)

	if err = v.votingService.DeleteVoting(c.Request.Context(), &DeleteVotingRequest{
		ID: id,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "voting deleted"})
}

type MakeChoiceRequest struct {
	InvarianceID uuid.UUID `json:"invarianceId"`
}

func (v *VotingHandler) MakeChoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	fmt.Println("id: ", id)

	var request MakeChoiceRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = v.votingService.MakeChoice(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "choice made"})
}
