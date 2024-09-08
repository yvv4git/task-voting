/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yvv4git/task-voting/internal/application"
	"github.com/yvv4git/task-voting/internal/infrastructure"
)

// votingCmd represents the voting command
var votingCmd = &cobra.Command{
	Use:   "voting",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := infrastructure.NewDefaultLogger()

		var config infrastructure.Config
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Error("unmarshalling config", err)
			return
		}

		appVoting := application.NewVoting(log, config)
		if err := appVoting.Start(); err != nil {
			log.Error("start application", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(votingCmd)
}
