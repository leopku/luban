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
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "github.com/src-d/go-mysql-server"
  "github.com/src-d/go-mysql-server/auth"
  "github.com/src-d/go-mysql-server/server"
  fakesql "github.com/src-d/go-mysql-server/sql"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
  Use:   "mysql",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  Run: runMysql,
}

func init() {
  rootCmd.AddCommand(mysqlCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // mysqlCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // mysqlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runMysql(cmd *cobra.Command, args []string) {
  driver := sqle.NewDefault()
  driver.AddDatabase(utils.CreateTestDatabase())
  driver.AddDatabase(fakesql.NewInformationSchemaDatabase(driver.Catalog))

  config := server.Config{
    Protocol: "tcp",
    Address:  "127.0.0.1:3366",
    Auth:     auth.NewNativeSingle("user", "pass", auth.AllPermissions),
  }

  s, err := server.NewDefaultServer(config, driver)
  if err != nil {
    log.Fatal().Err(err).Msg("")
  }

  fmt.Printf("Mysql now listening on: '%s'\n", config.Address)
  s.Start()
}
