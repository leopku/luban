package utils

import (
  "database/sql"
  // "github.com/rs/zerolog/log"
)

func ProviderDB() *sql.DB {
  return NewDB()
}
