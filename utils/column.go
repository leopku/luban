package utils

import (
  "database/sql"
  "io/ioutil"
  "strings"

  strUtil "github.com/agrison/go-commons-lang/stringUtils"
  "github.com/dave/jennifer/jen"
  // "github.com/huandu/xstrings"
  // "github.com/francoispqt/gojay"
  "github.com/iancoleman/strcase"
  "github.com/json-iterator/go"
  "github.com/rs/zerolog/log"
  "github.com/ulule/deepcopier"
  // "github.com/vitkovskii/insane-json"
)

type FieldMapping struct {
  SqlType        string `json:"sql_type" deepcopier:"field:SqlType"`
  GoType         string `json:"go_type" deepcopier:"field:GoType"`
  JsonType       string `json:"json_type" deepcopier:"field:JsonType"`
  ProtobufType   string `json:"protobuf_type" deepcopier:"field:ProtobufType"`
  GureguType     string `json:"guregu_type" deepcopier:"field:GureguType"`
  GoNullableType string `json:"go_nullable_type" deepcopier:"field:GoNullableType"`
  SwaggerType    string `json:"swagger_type" deepcopier:"field:SwaggerType"`
  Size           int    `json:"size,omitempty"`
  Custom         string `json:"custom,omitempty"`
}

type MappingFile struct {
  Mappings []*FieldMapping `json:"mappings"`
}

func NewMappingFromJSON(filename string) (*MappingFile, error) {
  buf, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  mappingFile := &MappingFile{}
  if err := jsoniter.Unmarshal(buf, mappingFile); err != nil {
    return nil, err
  }
  // if err := gojay.Unmarshal(buf, mappingFile); err != nil {
  /*if err := gojay.UnmarshalJSONObject(buf, mappingFile); err != nil {
    return nil, err
  }*/
  // root, err := insaneJSON.DecodeBytes(buf)
  // defer insaneJSON.Release(root)

  return mappingFile, nil
}

type ColumnMeta struct {
  FieldMapping
  Column *sql.ColumnType
}

func (this *ColumnMeta) GetName() string {
  return this.Column.Name()
}

func (this *ColumnMeta) GetGoName() string {
  ret := strcase.ToCamel(this.GetName())
  return ret
}

func (this *ColumnMeta) GetSqlType() string {
  return this.Column.DatabaseTypeName()
}

func (this *ColumnMeta) GetJsonName() string {
  return strcase.ToLowerCamel(this.GetName())
}

func (this *ColumnMeta) ParseAllTypes(maps []*FieldMapping) {
  for i, fMap := range maps {
    found := false
    if strings.ToLower(fMap.SqlType) == strings.ToLower(this.GetSqlType()) {
      deepcopier.Copy(fMap).To(this)
      log.Trace().Str("sql type", this.SqlType).Str("go type", this.GoType).Str("json type", this.JsonType).Msg("after copying")

      found = true
    }

    if found {
      break
    }

    if !found && i == len(maps) {
      this.SqlType = this.GetSqlType()
      this.GoType = "interface{}"
    }
  }
}

// type JenniferFn func(state *jen.Statement) *jen.Statement {}

func (this *ColumnMeta) GenerateGo(group *jen.Group) {
  if strUtil.IsBlank(this.GoType) {
    return
  }
  s := group.Id(this.GetGoName())
  nullable, ok := this.Column.Nullable()
  // jsonTagMap := map[string]string{"json": this.GetJsonName()}
  jsonTag := this.GetJsonName()
  if ok && nullable {
    jsonTag = jsonTag + ",omitempty"
  }

  switch this.GoType {
  case "string":
    if !ok || !nullable {
      s.String()
    } else {
      s.Qual("github.com/volatiletech/null", "String")
    }
  case "bool":
    if !ok || !nullable {
      s.Bool()
    } else {
      s.Qual("github.com/volatiletech/null", "Bool")
    }
  case "int32":
    if !ok || !nullable {
      s.Int32()
    } else {
      s.Qual("github.com/volatiletech/null", "Int32")
    }
  case "int64":
    if !ok || !nullable {
      s.Int64()
    } else {
      s.Qual("github.com/volatiletech/null", "Int64")
    }
  case "time.Time":
    if !ok || !nullable {
      s.Qual("time", "Time")
    } else {
      s.Qual("github.com/volatiletech/null", "Time")
    }
  case "float64":
    if !ok || !nullable {
      s.Float64()
    } else {
      s.Qual("github.com/volatiletech/null", "Float64")
    }
  case "float32":
    if !ok || !nullable {
      s.Float32()
    } else {
      s.Qual("github.com/volatiletech/null", "Float32")
    }
  case "[]byte":
    if !ok || !nullable {
      s.Index().Byte()
    } else {
      s.Index().Qual("github.com/volatiletech/null", "Byte")
    }
  case "uint32":
    if !ok || !nullable {
      s.Uint32()
    } else {
      s.Qual("github.com/volatiletech/null", "Uint32")
    }
  case "uint64":
    if !ok || !nullable {
      s.Uint64()
    } else {
      s.Qual("github.com/volatiletech/null", "Uint64")
    }
  case "interface{}":
  default:
    s.Interface()
  }

  s.Tag(map[string]string{"json": jsonTag, "db": this.GetName()})

}
