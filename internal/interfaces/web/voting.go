package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yvv4git/task-voting/internal/infrastructure"
)

type contextKey string

const userIDKey contextKey = "userID"

type VotingService interface {
	List(ctx context.Context, r *ListVotingRequest) (*ListVotingResponse, error)
	CreateVoting(ctx context.Context, r *CreateVotingRequest) (*CreateVotingResponse, error)
	UpdateVoting(ctx context.Context, r *UpdateVotingRequest) error
	DeleteVoting(ctx context.Context, r *DeleteVotingRequest) error
	MakeChoice(ctx context.Context, r *MakeChoiceRequest) error
}

type AuthService interface {
	CheckLoginPassword(ctx context.Context, login string, password string) error
	UserIDByLoginPassword(ctx context.Context, username, password string) (uuid.UUID, error)
}

type SubscriptionProcessor interface {
	AddClient(client infrastructure.ClientConn)
}

type VotingHandler struct {
	votingService VotingService
	authService   AuthService
	subscription  SubscriptionProcessor
}

func NewVotingHandler(votingService VotingService, authService AuthService, subscription SubscriptionProcessor) *VotingHandler {
	return &VotingHandler{
		votingService: votingService,
		authService:   authService,
		subscription:  subscription,
	}
}

func (v *VotingHandler) RegisterHandlers(router *gin.Engine) {
	authMiddleware := func(c *gin.Context) {
		login, password, err := infrastructure.ExtractBasicAuthValid(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check login and password with auth service
		if err := v.authService.CheckLoginPassword(c.Request.Context(), login, password); err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Request userID from auth service
		userID, err := v.authService.UserIDByLoginPassword(c.Request.Context(), login, password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Store the userID in the context
		ctx := context.WithValue(c.Request.Context(), userIDKey, userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

	votingGroup := router.Group("/voting")
	votingGroup.Use(authMiddleware)
	{
		votingGroup.GET("", v.ListVoting)
		votingGroup.POST("", v.CreateVoting)
		votingGroup.PUT("/:id", v.UpdateVoting)
		votingGroup.DELETE("/:id", v.DeleteVoting)
		votingGroup.POST("/choice/:id", v.MakeChoice)
		votingGroup.GET("/subscribe", v.Subscribe)
	}
}

type (
	ListVotingRequest struct {
		Limit  int64 `form:"limit"`
		Offset int64 `form:"offset"`
	}

	InvarianceScore struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Score int64     `json:"score"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := CreateVotingResponse{
		ID: result.ID,
	}

	c.JSON(http.StatusCreated, response)
}

type UpdateVotingRequest struct {
	ID          uuid.UUID  `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	StartAt     *time.Time `json:"startAt"`
	EndAt       *time.Time `json:"endAt"`
	Invariance  []string   `json:"invariance"`
}

func (v *VotingHandler) UpdateVoting(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var request UpdateVotingRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = v.votingService.UpdateVoting(c.Request.Context(), &UpdateVotingRequest{
		ID:          id,
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

	if err = v.votingService.DeleteVoting(c.Request.Context(), &DeleteVotingRequest{
		ID: id,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "voting deleted"})
}

type MakeChoiceRequest struct {
	InvarianceID uuid.UUID
	UserID       uuid.UUID
}

func (v *VotingHandler) MakeChoice(c *gin.Context) {
	idStr := c.Param("id")
	invarianceID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invariance id"})
		return
	}

	// Extract user id from context
	ctx := c.Request.Context()
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err = v.votingService.MakeChoice(c.Request.Context(), &MakeChoiceRequest{
		InvarianceID: invarianceID,
		UserID:       userID,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully voted"})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (v *VotingHandler) Subscribe(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade to websocket"})
		return
	}

	v.subscription.AddClient(ws)
}
