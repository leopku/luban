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
  "os"

  "github.com/leopku/luban/utils"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var (
  cfgFile string
  vip     *viper.Viper
  v       bool
  vv      bool
  vvv     bool
  vvvv    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "luban",
  Short: "A brief description of your application",
  Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //  Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
  // cobra.OnInitialize(initLog)
  utils.InitConfig(
    utils.WithV(v),
    utils.WithVV(vv),
    utils.WithV(vvv),
    utils.WithVVVV(vvvv),
  )

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.luban.yaml)")
  rootCmd.PersistentFlags().BoolVar(&v, "v", false, "show detailed output")
  rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "show more detailed output")
  rootCmd.PersistentFlags().BoolVar(&vvv, "vvv", false, "show more and more detailed output")
  rootCmd.PersistentFlags().BoolVar(&vvvv, "vvvv", false, "show most detailed output")

  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
  // rootCmd.Flags().BoolP("v", "", false, "show detailed output")
  // rootCmd.Flags().BoolP("vv", "", false, "show more detailed output")
  // rootCmd.Flags().BoolP("vvv", "", false, "show most detailed output")
  // rootCmd.Flags().BoolP("vvvv", "", false, "show most detailed output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name "." (without extension).
    viper.AddConfigPath(fmt.Sprintf("%s/.config", home))
    viper.AddConfigPath(".")
    viper.SetConfigName("luban")
  }

  viper.SetEnvPrefix("luban")
  viper.AutomaticEnv() // read in environment variables that match

  viper.SetDefault("database.adapter", "postgres")
  viper.SetDefault("database.host", "localhost")
  viper.SetDefault("database.port", 5432)
  if viper.GetString("database.adapter") == "mysql" {
    viper.SetDefault("database.port", "3306")
  }
  viper.SetDefault("database.encoding", "utf8mb4")
  viper.SetDefault("generation.model.output", "./models")

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
  vip = viper.GetViper()
}

func initLog() {
  level := zerolog.ErrorLevel
  if v {
    level = zerolog.WarnLevel
  }
  if vv {
    level = zerolog.InfoLevel
  }
  if vvv {
    level = zerolog.DebugLevel
  }
  if vvvv {
    level = zerolog.TraceLevel
  }
  zerolog.SetGlobalLevel(level)

  log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
