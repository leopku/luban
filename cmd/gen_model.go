package cmd

import (
  // "fmt"

  // "github.com/iancoleman/strcase"
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

  // output := viper.GetString("generation.model.output")

  tables, err := utils.GetAllTableMeta(db, vip)
  if err != nil {
    return
  }

  for _, table := range tables {
    output := viper.GetString("generation.output")
    //usecasePath := fmt.Sprintf("%s/usecase/%s", output, strcase.ToKebab(table.GetNameWithoutPrefix()))
    modulePath := table.GetModulePath(output)
    if err := utils.CreateDirectory(modulePath); err != nil {
      log.Fatal().Err(err).Msg("")
    }

    // modelFile := fmt.Sprintf("%s/%s_model.go", usecasePath, table.GetGoFileName())
    err := table.SaveModel(modulePath)
    if err != nil {
      log.Error().Err(err).Str("table model", table.Name).Msg("generating failed")
    } else {
      count++
    }
  }

  log.Log().Int("count", count).Msg("Models generating successfully")
}
