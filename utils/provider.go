package utils

import (
  "database/sql"
  "io"

  "github.com/spf13/viper"
)

func ProviderConfig(in io.Reader, cType string) *viper.Viper {
  return NewConfigFromReader(in, cType)
}

func ProviderDB(vip *viper.Viper) *sql.DB {
  return NewDB(vip)
}
