// +build wireinject

package utils

import (
  "database/sql"
  "github.com/google/wire"
)

func BuildDB() *sql.DB {
  wire.Build(ProviderDB)
  return nil
}
