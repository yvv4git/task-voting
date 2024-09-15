package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/yvv4git/task-voting/internal/domain/repository"
	"github.com/yvv4git/task-voting/internal/domain/service"
	"github.com/yvv4git/task-voting/internal/infrastructure"
	"github.com/yvv4git/task-voting/internal/interfaces/web"
)

type VotingService interface {
}

type Voting struct {
	application
	cfg infrastructure.Config
}

func NewVoting(log *slog.Logger, cfg infrastructure.Config) *Voting {
	v := &Voting{
		application: application{
			log: log,
		},
		cfg: cfg,
	}
	v.app = v
	return v
}

func (v *Voting) start(ctx context.Context) error {
	v.log.Info("Starting VotingApplication")
	defer v.log.Info("Shutting down VotingApplication")

	// Init db
	db, err := infrastructure.NewPostgresDB(ctx, v.cfg.VotingApp.DataBase)
	if err != nil {
		return fmt.Errorf("init db connection: %w", err)
	}

	// Init repo & service
	votingRepo := repository.NewVoting(db)
	votingService := service.NewVoting(votingRepo)
	authService := infrastructure.NewAuthStub()

	// Init web interface
	webConfig := v.cfg.VotingApp.WebAPI
	r := gin.Default()
	webHandler := web.NewVotingHandler(votingService, authService)
	webHandler.RegisterHandlers(r)
	webSrv := infrastructure.NewWebServer(v.log, r, fmt.Sprintf("%s:%d", webConfig.Host, webConfig.Port))
	if err = webSrv.Run(ctx); err != nil {
		return err
	}

	return nil
}
