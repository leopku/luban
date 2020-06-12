package cmd

import (
  "github.com/leopku/luban/utils"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var genModelCmd = &cobra.Command{
  Use:   "model",
  Short: "generate models",
  Long:  `generate models`,
  Run:   runGenModel,
}

func init() {
  genCmd.AddCommand(genModelCmd)
}

func runGenModel(cmd *cobra.Command, args []string) {
  count := 0
  db := utils.NewDB(vip)
  if db == nil {
    log.Fatal().Msg("db init failed")
  }

  var err error
  defer func() {
    if err != nil {
      log.Fatal().Err(err).Msg("")
    }
  }()

  output := viper.GetString("generation.model.output")
  if err := utils.CreateDirectory(output); err != nil {
    log.Fatal().Err(err).Msg("")
  }

  tables, err := utils.GetAllTableMeta(db, vip)
  if err != nil {
    return
  }

  for _, table := range tables {
    if err := table.SaveToGo(output); err != nil {
      log.Error().Err(err).Str("table", table.Name).Msg("generating failed")
    } else {
      count++
    }
  }

  log.Log().Int("count", count).Msg("Models generating successfully")
}
