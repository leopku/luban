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
  // "errors"
  "fmt"
  // "path"
  "github.com/leopku/luban/generated/module/ad"
  "github.com/leopku/luban/generated/module/coupon"
  "github.com/leopku/luban/utils"
  // "github.com/dave/jennifer/jen"
  // // "github.com/iancoleman/strcase"
  // // "github.com/novalagung/gubrak/v2"
  // strUtil "github.com/agrison/go-commons-lang/stringUtils"
  // "github.com/iancoleman/strcase"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  // "github.com/spf13/viper"
  // "gopkg.in/src-d/go-parse-utils.v1"
  "github.com/jmoiron/sqlx"
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

    dbx := sqlx.NewDb(db, "mysql")
    // adRepo := &ad.AdRepository{db: dbx}
    adRepo := ad.NewAdRepository(dbx)
    adSlice, err := adRepo.GetAll()
    if err != nil {
      log.Fatal().Err(err).Msg("")
    }
    log.Log().Interface("all ad", adSlice).Msg("")

    ad, err := adRepo.GetById(1)
    log.Log().Interface("ad by id", ad).Msg("")

    couponRepo := coupon.NewCouponRepository(dbx)
    couponSlice, err := couponRepo.GetAll(coupon.WithWhere("`limit` = %d", 1))
    // couponSlice, err := couponRepo.GetAll(coupon.WithWhere("`limit` = %d AND type = %d", 1, 1))
    if err != nil {
      log.Fatal().Err(err).Msg("")
    }
    log.Log().Interface("all coupon", couponSlice).Msg("")
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
