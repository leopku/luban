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
  // f := jen.NewFile(GetModelPackageName(this.OutputPath) + "/" + this.GetNameWithoutPrefix())
  f := jen.NewFile(this.GetNameWithoutPrefix())
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
  // f := jen.NewFile(GetModelPackageName(this.OutputPath) + "/" + this.GetNameWithoutPrefix())
  f := jen.NewFile(this.GetNameWithoutPrefix())
  f.ImportAlias("github.com/Masterminds/squirrel", "sql")

  // type Option
  f.Line().Comment("option")
  f.Type().Id("Option").Struct(
    jen.Id("Where").String(),
    jen.Id("Deleted").Bool(),
    jen.Id("Limit").Int(),
    jen.Id("Offset").Int(),
    jen.Id("OrderBy").String(),
    jen.Id("Order").String(),
  )
  // type OptionFunc
  f.Type().Id("OptionFunc").Func().
    Params(
      jen.Id("opts").Op("*").Id("Option"),
    )

  // var defaultOption
  f.Var().Id("defaultOption").Op("=").Id("Option").Values(jen.Dict{
    jen.Id("Where"):   jen.Lit("1=1"),
    jen.Id("Deleted"): jen.Lit(false),
    jen.Id("Limit"):   jen.Lit(10),
    jen.Id("Offset"):  jen.Lit(0),
    jen.Id("OrderBy"): jen.Lit(""),
    jen.Id("Order"):   jen.Lit("asc"),
  })

  // func WithLimit
  f.Func().Id("WithLimit").
    Params(jen.Id("limit").Int()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          jen.Id("opts").Dot("Limit").Op("=").Id("limit"),
        ),
      ),
    )

  // func WithOffset
  f.Func().Id("WithOffset").
    Params(jen.Id("offset").Int()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          jen.Id("opts").Dot("Offset").Op("=").Id("offset"),
        ),
      ),
    )

  // func WithOrderBy
  f.Func().Id("WithOrderBy").
    Params(jen.Id("orderBy").String()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          jen.Id("opts").Dot("OrderBy").Op("=").Id("orderBy"),
        ),
      ),
    )

  // func WithOrder
  f.Func().Id("WithOrder").
    Params(jen.Id("order").String()).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          jen.Id("opts").Dot("Order").Op("=").Id("order"),
        ),
      ),
    )

  // func WithWhere
  f.Func().Id("WithWhere").
    Params(
      jen.List(
        jen.Id("format").String(),
        jen.Id("a").Op("...").Interface(),
      ),
    ).
    Params(jen.Id("OptionFunc")).
    Block(
      jen.Return(
        jen.Func().Params(jen.Id("opts").Op("*").Id("Option")).Block(
          jen.Id("opts").Dot("Where").Op("=").Qual("fmt", "Sprintf").Params(
            jen.List(
              jen.Lit("%s AND 1=1"),
              jen.Qual("fmt", "Sprintf").Params(jen.List(jen.Id("format"), jen.Id("a").Op("..."))),
            ),
          ),
        ),
      ),
    )

  // interface IxxxRepostiry
  f.Line().Comment(fmt.Sprintf("%s respostiory interface", this.GetModelName()))
  f.Type().Id(fmt.Sprintf("I%sRepository", this.GetModelName())).Interface(
    jen.Id("GetAll").
      Params(
        jen.Id("opts").Op("...").Id("OptionFunc"),
      ).
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

  repositoryName := fmt.Sprintf("%sRepository", this.GetModelName())
  // repository xxxRepository
  f.Line().Comment(fmt.Sprintf("%s repository", this.GetModelName()))
  f.Type().Id(repositoryName).Struct(
    jen.Id("db").Op("*").Qual("github.com/jmoiron/sqlx", "DB"),
  )

  // func NewxxxRepository
  f.Line().Comment(fmt.Sprintf("New%sRepository", this.GetModelName()))
  f.Func().Id(fmt.Sprintf("New%s", repositoryName)).
    Params(jen.Id("db").Op("*").Qual("github.com/jmoiron/sqlx", "DB")).
    Params(jen.Op("*").Id(repositoryName)).
    Block(jen.Return(
      jen.Op("&").
        Id(repositoryName).
        Values(jen.Dict{jen.Id("db"): jen.Id("db")}),
    ))

  // func GetAll
  f.Line().Comment(fmt.Sprintf("get all %s", this.GetModelName()))
  f.Func().
    Params(jen.Id("this").Op("*").Id(repositoryName)).
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
      jen.Id(fmt.Sprintf("%sSlice", this.GetNameWithoutPrefix())).Op(":=").Index().Op("*").Id(this.GetModelName()).Block(),
      jen.Id("sql").Op(":=").Lit(fmt.Sprintf("SELECT * FROM `%s`", this.Name)),
      jen.Id("sql").Op("=").Qual("fmt", "Sprintf").
        Params(jen.List(
          jen.Lit("%s WHERE %s"),
          jen.Id("sql"),
          jen.Id("options").Dot("Where"),
        )),
      jen.Id("sql").Op("=").Qual("fmt", "Sprintf").
        Params(jen.List(
          jen.Lit("%s LIMIT %d OFFSET %d"),
          jen.Id("sql"),
          jen.Id("options").Dot("Limit"),
          jen.Id("options").Dot("Offset"),
        )),
      jen.If(
        jen.Op("!").Qual("github.com/agrison/go-commons-lang/stringUtils", "IsBlank").
          Params(jen.Id("options").Dot("OrderBy")).
          Block(
            jen.Id("sql").Op("=").Qual("fmt", "Sprintf").
              Params(jen.List(
                jen.Lit("%s ORDER BY %s %s"),
                jen.Id("sql"),
                jen.Id("options").Dot("OrderBy"),
                jen.Id("options").Dot("Order"),
              )),
          ),
      ),
      jen.Err().Op(":=").Id("this").Dot("db").Dot("Select").Call(
        jen.Op("&").Id(fmt.Sprintf("%sSlice", this.GetNameWithoutPrefix())),
        jen.Id("sql"),
      ),
      jen.If(
        jen.Err().Op("!=").Id("nil").Block(
          jen.Return(jen.Nil(), jen.Err()),
        ),
      ),
      jen.Return(
        jen.Id(fmt.Sprintf("%sSlice", this.GetNameWithoutPrefix())),
        jen.Nil()))

  // func GetById
  f.Line().Comment(fmt.Sprintf("get one %s by id", this.GetModelName()))
  f.Func().
    Params(jen.Id("this").Op("*").Id(repositoryName)).
    Id("GetById").
    Params(jen.Id("id").Int()).
    Params(jen.List(
      jen.Op("*").Id(this.GetModelName()),
      jen.Error(),
    )).
    Block(
      jen.Id(this.GetNameWithoutPrefix()).Op(":=").Op("&").Id(this.GetModelName()).Block(),
      // TODO: add deleted_at / deleted in where clause
      jen.Id("sb").Op(":=").
        Qual("github.com/Masterminds/squirrel", "Select").Call(jen.Lit("*")).
        Dot("From").Call(jen.Lit(fmt.Sprintf("%s", this.Name))).
        Dot("Where").Call(jen.Lit("id = ?"), jen.Id("id")),
      jen.List(jen.Id("query"), jen.Id("args"), jen.Err()).Op(":=").Id("sb").Dot("ToSql").Call(),
      jen.If(
        jen.Err().Op("!=").Nil().Block(
          // jen.Qual("fmt", "Println").Params(jen.Id("query")),
          jen.Return(jen.Nil(), jen.Err()),
        ),
      ),
      jen.If(
        jen.Err().Op(":=").Id("this").Dot("db").Dot("Get").Call(
          jen.Id(this.GetNameWithoutPrefix()),
          jen.Id("query"),
          jen.Id("args").Op("..."),
        ),
        jen.Err().Op("!=").Nil(),
      ).
        Block(jen.Return(jen.Nil(), jen.Err())),
      jen.Return(jen.Id(this.GetNameWithoutPrefix()), jen.Nil()),
    )

  return f
}

func (this *TableMeta) SaveRepo(modulePath string) error {
  f := this.BuildRepo()
  return f.Save(this.GetFileName(modulePath, "repo"))
}
