/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
  "fmt"

  "github.com/leopku/luban/utils"

  // "github.com/dave/jennifer/jen"
  // "github.com/iancoleman/strcase"
  // "github.com/novalagung/gubrak/v2"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
  Use:   "demo",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("demo called")
    fmt.Println("ConfigFileUsed: ", viper.ConfigFileUsed())
    fmt.Println(viper.GetString("database.aa"))
    log.Trace().Interface("database", viper.Get("database")).Msg("")
    // fmt.Println(zerolog.GetLevel().String())
    db := utils.BuildDB()
    if db == nil {
      log.Fatal().Msg("db init failed")
    }

    var err error
    defer func() {
      if err != nil {
        log.Fatal().Err(err).Msg("")
      }
    }()

    tables, err := utils.GetAllTableMeta(db)
    if err != nil {
      return
    }
    t := tables[0]
    log.Trace().Str("table name", t.Name).Str("model name", t.GetModelName()).Str("go filename", t.GetGoFileName()+".go").Msg("")
    cols, err := t.GetAllColumnMeta()
    for _, col := range cols {
      // log.Trace().Str("column name", col.GetName()).Str("sql type", col.SqlType).Str("go type", col.GoType).Str("json type", col.JsonType).Msg("")
      log.Debug().Str("column name", col.GetName()).Str("sql type", col.SqlType).Str("go type", col.GoType).Str("json type", col.JsonType).Msg("")
    }

    mapping, err := utils.NewFromJSON("./templates/mapping.json")
    if err != nil {
      return
    }
    log.Trace().Interface("mapping", mapping).Msg("")

  },
}

func init() {
  rootCmd.AddCommand(demoCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // demoCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // demoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
