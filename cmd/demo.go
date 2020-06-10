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
  "strings"

  "github.com/leopku/luban/utils"

  "github.com/jimsmart/schema"
  "github.com/novalagung/gubrak/v2"
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
    tableSlice, err := schema.TableNames(db)
    utils.IfErrCallback(err, func() interface{} {
      log.Fatal().Err(err).Msg("")
      return nil
    })
    log.Trace().Interface("tables", tableSlice).Msg("")

    // excludes, ok := viper.Get("generation.exclude").([]string)
    // excludes := viper.Get("generation.exclude")
    excludes := viper.GetStringSlice("generation.exclude")
    // if !ok {
    //   log.Error().Msg("exclude param error, ignore while processing")
    // }
    log.Debug().Interface("exclude", excludes).Msg("")
    for _, v := range excludes {
      log.Debug().Str("exclude", v).Msg("")
    }
    retInter := gubrak.From(tableSlice).
      Intersection(excludes).
      Result()
    log.Debug().Interface("intersection", retInter).Msg("")
    /*retExclude := gubrak.From(tableSlice).
      Exclude(excludes).
      Result()*/
    chainableExclude := gubrak.From(tableSlice)
    for _, v := range excludes {
      chainableExclude.Exclude(v)
    }
    retExclude := chainableExclude.Result()
    log.Debug().Interface("exclude", retExclude).Msg("")
    for _, table := range retExclude.([]string) {
      log.Debug().Str("prefix stripped", strings.TrimPrefix(table, viper.GetString("generation.prefix"))).Msg("")
    }
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
