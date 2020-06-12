// +build wireinject

package utils

import (
  "database/sql"
  "io"

  "github.com/google/wire"
  "github.com/spf13/viper"
)

func BuildConfig(in io.Reader, cType string) *viper.Viper {
  wire.Build(ProviderConfig)
  return nil
}

func BuildDB(in io.Reader, cType string) *sql.DB {
  wire.Build(ProviderConfig, ProviderDB)
  return nil
}
