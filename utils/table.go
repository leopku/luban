package utils

import (
  "database/sql"
  // "errors"
  "fmt"
  "strings"

  // strUtil "github.com/agrison/go-commons-lang/stringUtils"
  "github.com/dave/jennifer/jen"
  // "github.com/francoispqt/gojay"
  "github.com/gertd/go-pluralize"
  "github.com/huandu/xstrings"
  "github.com/iancoleman/strcase"
  "github.com/jimsmart/schema"
  "github.com/rs/zerolog/log"
  // "github.com/volatiletech/null"
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

func (this *TableMeta) GetModuleName() string {
  return strcase.ToKebab(this.GetNameWithoutPrefix())
}

func (this *TableMeta) GetModulePath(base string) string {
  return fmt.Sprintf("%s/module/%s", base, this.GetModuleName())
}

func (this *TableMeta) GetFileName(modulePath string, genre string) string {
  return fmt.Sprintf("%s/%s_%s.go", modulePath, this.GetGoFileName(), genre)
}

func (this *TableMeta) GetAllColumnMeta() ([]*ColumnMeta, error) {
  cols, err := schema.Table(this.db, this.Name)
  if err != nil {
    return nil, err
  }

  ret := []*ColumnMeta{}
  mapping, err := NewMappingFromJSON("./templates/mapping.json")
  // if errors.Is(err, gojay.InvalidUnmarshalError) {
  // if err == gojay.InvalidUnmarshalError {
  //   log.Error().Msg("InvalidUnmarshalError")
  // }
  if err != nil {
    return nil, err
  }
  log.Trace().Interface("mapping", mapping).Msg("")

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
    return f
  }

  f.Line().Comment(fmt.Sprintf("model %s", this.GetModelName()))
  f.Type().Id(this.GetModelName()).StructFunc(func(g *jen.Group) {
    for _, col := range columns {
      log.Trace().Str("col name", col.GetGoName()).Str("col sql type", col.SqlType).Msg("")
      col.GenerateGo(g)
    }
  })

  return f
}

func (this *TableMeta) SaveModel(modulePath string) error {
  f := this.BuildModel()
  return f.Save(this.GetFileName(modulePath, "model"))
}

func (this *TableMeta) BuildRepo() *jen.File {
  f := jen.NewFile(GetModelPackageName(this.OutputPath))

  f.Line().Comment("option")
  f.Type().Id("Option").Struct(
    jen.Id("Limit").Int(),
    jen.Id("Offset").Int(),
    jen.Id("Order").String(),
  )
  f.Type().Id("OptionFunc").Func().
    Params(
      jen.Id("opts").Op("*").Id("Option"),
    )
  f.Var().Id("defaultOption").Op("=").Id("Option").Values(jen.Dict{
    jen.Id("Limit"):  jen.Lit(10),
    jen.Id("Offset"): jen.Lit(0),
    jen.Id("Order"):  jen.Lit("asc"),
  })
  f.Func().Id("WithLimit").
    Params(jen.Id("limit").Int()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          // jen.Qual("opts", "Limit").Op("=").Id("limit"),
          jen.Id("opts").Dot("Limit").Op("=").Id("limit"),
        ),
      ),
    )
  f.Func().Id("WithOffset").
    Params(jen.Id("offset").Int()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          // jen.Qual("opts", "Offset").Op("=").Id("offset"),
          jen.Id("opts").Dot("Offset").Op("=").Id("offset"),
        ),
      ),
    )
  f.Func().Id("WithOrder").
    Params(jen.Id("order").Int()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          // jen.Qual("opts", "Order").Op("=").Id("order"),
          jen.Id("opts").Dot("Order").Op("=").Id("order"),
        ),
      ),
    )

  f.Line().Comment(fmt.Sprintf("%s respostiory interface", this.GetModelName()))
  f.Type().Id(fmt.Sprintf("I%sRepository", this.GetModelName())).Interface(
    jen.Id("GetAll").
      Params().
      Params(
        jen.Index().Op("*").Id(this.GetModelName()),
        jen.Error(),
      ),
    jen.Id("GetById").
      Params(jen.Id("id").Int32()).
      Params(
        jen.Op("*").Id(this.GetModelName()),
        jen.Error(),
      ),
    jen.Id("Create").
      Params(jen.Id(strcase.ToLowerCamel(this.GetModelName())).Op("*").Id(this.GetModelName())).
      Params(jen.Error()),
    jen.Id("Update").
      Params(
        jen.Id("id").Int32(),
        // jen.Id(strcase.ToLowerCamel(this.GetModelName()))).Op("*").Id(this.GetModelName()).
        jen.Id(strcase.ToLowerCamel(this.GetModelName())).Op("*").Id(this.GetModelName())).
      Params(jen.Error()),
    jen.Id("Delete").
      Params(jen.Id("id").Int32()).
      Params(jen.Error()),
  )

  f.Line().Comment(fmt.Sprintf("%s repository", this.GetModelName()))
  f.Type().Id(fmt.Sprintf("%sRepository", this.GetModelName())).Struct(
    jen.Id("db").Op("*").Qual("database/sql", "DB"),
  )

  f.Line().Comment(fmt.Sprintf("%s repository", this.GetModelName()))
  f.Func().
    Params(jen.Id("this").Op("*").Id(fmt.Sprintf("%sRepository", this.GetModelName()))).
    Id("GetAll").
    Params(
      jen.Id("opts").Op("...").Id("OptionFunc"),
    ).
    Params(
      jen.Index().Op("*").Id(this.GetModelName()),
      jen.Error()).
    Block(
      jen.Id("options").Op(":=").Id("defaultOption"),
      jen.For(
        jen.List(jen.Id("_"), jen.Id("opt")).Op(":=").Range().Id("opts").Block(
          jen.Id("opt").Params(
            jen.Op("&").Id("options"),
          ),
        ),
      ),
      jen.Id("sql").Op(":=").Lit(fmt.Sprintf("SELECT * FROM `%s`", this.Name)),
      jen.Return(
        jen.Nil(),
        jen.Nil()))

  return f
}

func (this *TableMeta) SaveRepo(modulePath string) error {
  f := this.BuildRepo()
  return f.Save(this.GetFileName(modulePath, "repo"))
}
