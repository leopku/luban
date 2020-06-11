package utils

import (
  "database/sql"
  "fmt"

  // "github.com/CloudyKit/jet"
  _ "github.com/go-sql-driver/mysql"
  "github.com/jimsmart/schema"
  "github.com/rs/zerolog/log"
  "github.com/spf13/viper"
)

func NewDB() *sql.DB {
  connStr := ""
  adapter := viper.GetString("database.adapter")
  log.Trace().Str("adapter", adapter).Msg("")

  switch adapter {
  case "mysql":
    connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.database"), viper.GetString("database.encoding"))
  case "postgres":
    connStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.database"), viper.GetString("database.username"), viper.GetString("database.password"))
  }
  log.Trace().Str("connection string", connStr).Msg("")
  db, err := sql.Open(adapter, connStr)
  if err != nil {
    return nil
  }
  return db
}

func GetAllTableMeta(db *sql.DB) ([]*TableMeta, error) {
  ret := []*TableMeta{}
  tableSlice, err := schema.TableNames(db)
  if err != nil {
    return nil, err
  }
  prefix := viper.GetString("generation.prefix")
  for _, tName := range tableSlice {
    found := false
    for _, exclude := range viper.GetStringSlice("generation.exclude") {
      if tName == exclude {
        found = true
        break
      }
    }
    if !found {
      meta := &TableMeta{
        db:     db,
        Name:   tName,
        Prefix: prefix,
      }
      ret = append(ret, meta)
    }
  }
  return ret, nil
}
