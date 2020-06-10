package utils

import (
  "database/sql"
  "fmt"
  "strings"

  // "github.com/CloudyKit/jet"
  "github.com/emirpasic/gods/lists/arraylist"
  "github.com/gertd/go-pluralize"
  _ "github.com/go-sql-driver/mysql"
  "github.com/huandu/xstrings"
  "github.com/jimsmart/schema"
  "github.com/rs/zerolog/log"
  "github.com/spf13/viper"
  "github.com/ulule/deepcopier"
)

var plural = pluralize.NewClient()

type TableMeta struct {
  db     *sql.DB
  Name   string
  Prefix string
}

func (this *TableMeta) GetNameWithoutPrefix() string {
  return strings.TrimPrefix(this.Name, this.Prefix)
}

func (this *TableMeta) GetModelName() string {
  camelCase := xstrings.ToCamelCase(this.GetNameWithoutPrefix())
  return plural.Singular(camelCase)
}

func (this *TableMeta) GetGoFileName() string {
  snakeCase := xstrings.ToSnakeCase(this.GetNameWithoutPrefix())
  return plural.Singular(snakeCase)
}

func (this *TableMeta) GetAllColumnMetaA() ([]*sql.ColumnType, error) {
  return schema.Table(this.db, this.Name)
}

func (this *TableMeta) GetAllColumnMeta() ([]*ColumnMeta, error) {
  cols, err := schema.Table(this.db, this.Name)
  if err != nil {
    return nil, err
  }

  ret := []*ColumnMeta{}
  mapping, err := NewFromJSON("./templates/mapping.json")
  log.Trace().Interface("mapping", mapping).Msg("")
  if err != nil {
    return nil, err
  }

  for _, col := range cols {
    meta := &ColumnMeta{Column: col}
    meta.ParseAllTypes(mapping.Mappings)
    ret = append(ret, meta)
  }
  return ret, nil
}

type ColumnMeta struct {
  FieldMapping
  Column *sql.ColumnType
}

func (this *ColumnMeta) GetName() string {
  return this.Column.Name()
}

func (this *ColumnMeta) GetSqlType() string {
  return this.Column.DatabaseTypeName()
}

func (this *ColumnMeta) ParseAllTypes(maps []*FieldMapping) {
  for _, fMap := range maps {
    if strings.ToLower(fMap.SqlType) == strings.ToLower(this.GetSqlType()) {
      deepcopier.Copy(fMap).To(this)
      // this.SqlType = fMap.SqlType

      break
    }
  }
}

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

func GetAllTableNames(db *sql.DB) (tableNameList *arraylist.List, modelNameList *arraylist.List, modelFilenameList *arraylist.List, err error) {
  tableNameList = arraylist.New()
  modelNameList = arraylist.New()
  modelFilenameList = arraylist.New()

  tableSlice, err := schema.TableNames(db)
  if err != nil {
    return
  }

  prefix := viper.GetString("generation.prefix")
  plural := pluralize.NewClient()
  for _, tableName := range tableSlice {
    found := false
    for _, exclude := range viper.GetStringSlice("generation.exclude") {
      if tableName == exclude {
        found = true
        break
      }
    }

    if !found {
      tableNameList.Add(tableName)
      tableNameWithoutPrefix := strings.TrimPrefix(tableName, prefix)
      camelCase := xstrings.ToCamelCase(tableNameWithoutPrefix)
      snakeCase := xstrings.ToSnakeCase(tableNameWithoutPrefix)
      camelCaseSingular := plural.Singular(camelCase)
      snakeCaseSingular := plural.Singular(snakeCase)
      modelNameList.Add(camelCaseSingular)
      modelFilenameList.Add(snakeCaseSingular)
    }
  }

  return
}
