/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate a basic auth string",
	Long: `Generate a basic auth string from a login and password.

Example:
  go run main.go auth -l user -p password
  go run main.go auth -l user1 -p secret1
`,
	Run: func(cmd *cobra.Command, args []string) {
		login, err := cmd.Flags().GetString("login")
		if err != nil {
			fmt.Println(err)
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
			return
		}

		authString := login + ":" + password
		encodedAuthString := base64.StdEncoding.EncodeToString([]byte(authString))

		fmt.Println("Authorization: Basic", encodedAuthString)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Add flags for login and password
	authCmd.Flags().StringP("login", "l", "", "Login")
	authCmd.Flags().StringP("password", "p", "", "Password")
}
