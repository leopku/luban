package cmd

import (
  // "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  // "github.com/spf13/viper"
)

var genCmd = &cobra.Command{
  Use:   "gen",
  Short: "generate subcommands, using -h see all subcommands",
  Long:  `generate DDD layers like models, repositories, services`,
}

func init() {
  rootCmd.AddCommand(genCmd)
}
