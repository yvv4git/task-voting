/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yvv4git/task-voting/internal/application"
	"github.com/yvv4git/task-voting/internal/infrastructure"
)

// votingCmd represents the voting command
var votingCmd = &cobra.Command{
	Use:   "voting",
	Short: "Execute the application startup sequence",
	Long: `This command starts the voting service, initializing necessary components and beginning the application's startup sequence. 
	The service allows users to participate in voting processes, creating and managing voting sessions, casting votes, and retrieving results. 
	By executing this command, the application sets up the environment, establishes database connections, and initializes the voting service, ensuring its proper functioning.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := infrastructure.NewDefaultLogger()

		var config infrastructure.Config
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Error("unmarshalling config", slog.Any("error", err))
			return
		}

		appVoting := application.NewVoting(log, config)
		if err := appVoting.Start(); err != nil {
			log.Error("failed to start application", slog.Any("error", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(votingCmd)
}
