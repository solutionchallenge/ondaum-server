package main

import (
	"github.com/solutionchallenge/ondaum-server/cmd"
	"github.com/spf13/cobra"
)

// @title Ondaum API
// @version 1.0
// @description This is a API server for Ondaum
// @host ondaum.revimal.me
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	root := &cobra.Command{
		Use: "ondaum",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	root.AddCommand(cmd.NewConfigCommand())
	root.AddCommand(cmd.NewHttpCommand())

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
