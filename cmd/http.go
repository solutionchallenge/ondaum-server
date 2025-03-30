package cmd

import (
	"github.com/solutionchallenge/ondaum-server/internal/entrypoint/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/spf13/cobra"
)

func NewHttpCommand() *cobra.Command {
	httpCommand := &cobra.Command{
		Use: "http",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig := http.AppConfig{}
			cfgpath, err := cmd.Flags().GetString("config-path")
			if err != nil {
				panic(err)
			}
			cfgname, err := cmd.Flags().GetString("config-name")
			if err != nil {
				panic(err)
			}
			utils.LoadConfigTo(&appConfig, cfgname, cfgpath)
			http.Run(appConfig)
		},
	}
	httpCommand.Flags().StringP("config-path", "p", "./config", "config file path (default is './config')")
	httpCommand.Flags().StringP("config-name", "n", "production", "config file name (default is 'production')")
	return httpCommand
}
