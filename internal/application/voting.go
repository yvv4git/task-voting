package application

import (
	"context"
	"log/slog"

	"github.com/yvv4git/task-voting/internal/infrastructure"
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

	//v.log.Info("db.dbname", "value", v.cfg.VotingApp.DataBase.DBName)
	//v.log.Info("db.host", "value", v.cfg.VotingApp.DataBase.Host)
	//v.log.Info("db.port", "value", v.cfg.VotingApp.DataBase.Port)
	//v.log.Info("db.username", "value", v.cfg.VotingApp.DataBase.Username)
	//v.log.Info("db.password", "value", v.cfg.VotingApp.DataBase.Password)
	//
	//// Выводим параметры веб-API
	//v.log.Info("webapi.host", "value", v.cfg.VotingApp.WebAPI.Host)
	//v.log.Info("webapi.port", "value", v.cfg.VotingApp.WebAPI.Port)
	//v.log.Info("webapi.read_timeout", "value", v.cfg.VotingApp.WebAPI.ReadTimeout)
	//v.log.Info("webapi.shutdown_timeout", "value", v.cfg.VotingApp.WebAPI.ShutdownTimeout)

	// <-ctx.Done() check that signal works

	return nil
}
