package utils

import (
  "database/sql"
  "fmt"
  "strings"

  strUtil "github.com/agrison/go-commons-lang/stringUtils"
  "github.com/dave/jennifer/jen"
  "github.com/gertd/go-pluralize"
  "github.com/huandu/xstrings"
  "github.com/jimsmart/schema"
  "github.com/rs/zerolog/log"
  // "github.com/spf13/viper"
)

var plural = pluralize.NewClient()

type TableMeta struct {
  db         *sql.DB
  Name       string
  Prefix     string
  OutputPath string
  Columns    []*ColumnMeta
}

func (this *TableMeta) BuildName(name string) *TableMeta {
  this.Name = name
  return this
}

func (this *TableMeta) BuildPrefix(prefix string) *TableMeta {
  this.Prefix = prefix
  return this
}

func (this *TableMeta) GetNameWithoutPrefix() string {
  log.Trace().Str("prefix", this.Prefix).Msg("")
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

func (this *TableMeta) BuildModel() *jen.File {
  f := jen.NewFile(GetModelPackageName(this.OutputPath))
  // this.Columns := this.GetAllColumnMeta()
  columns, err := this.GetAllColumnMeta()
  if err != nil {
    log.Error().Err(err).Msg("")
    log.Warn().Msg("erorr caused this table only can generate empty struc.")
  }
  f.Type().Id(this.GetModelName()).StructFunc(func(g *jen.Group) {
    for _, col := range columns {
      log.Trace().Str("col name", col.GetGoName()).Msg("")
      if strUtil.IsBlank(col.GoType) {
        continue
      }
      s := g.Id(col.GetGoName())
      switch col.GoType {
      case "string":
        s.String()
      case "bool":
        s.Bool()
      case "int32":
        s.Int32()
      case "int64":
        s.Int64()
      case "time.Time":
        s.Qual("time", "Time")
      case "float64":
        s.Float64()
      case "float32":
        s.Float32()
      case "[]byte":
        s.Index().Byte()
      case "interface{}":
        s.Interface()
      case "uint32":
        s.Uint32()
      case "uint64":
        s.Uint64()
      default:
        s.Interface()
      }
      s.Tag(map[string]string{"json": col.GetJsonName()})
    }
  })

  return f
}

func (this *TableMeta) SaveToGo(path string) error {
  fullName := fmt.Sprintf("%s/%s.go", path, this.GetGoFileName())
  f := this.BuildModel()
  return f.Save(fullName)
}
