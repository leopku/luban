package utils

import (
  "database/sql"
  "io/ioutil"
  "strings"

  // strUtil "github.com/agrison/go-commons-lang/stringUtils"
  // "github.com/huandu/xstrings"
  "github.com/iancoleman/strcase"
  "github.com/json-iterator/go"
  "github.com/rs/zerolog/log"
  "github.com/ulule/deepcopier"
)

type FieldMapping struct {
  SqlType        string `json:"sql_type" deepcopier:"field:SqlType"`
  GoType         string `json:"go_type" deepcopier:"field:GoType"`
  JsonType       string `json:"json_type" deepcopier:"field:JsonType"`
  ProtobufType   string `json:"protobuf_type" deepcopier:"field:ProtobufType"`
  GureguType     string `json:"guregu_type" deepcopier:"field:GureguType"`
  GoNullableType string `json:"go_nullable_type" deepcopier:"field:GoNullableType"`
  SwaggerType    string `json:"swagger_type" deepcopier:"field:SwaggerType"`
}

type MappingFile struct {
  Mappings []*FieldMapping `json:"mappings"`
}

func NewFromJSON(filename string) (*MappingFile, error) {
  buf, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  // mappings := []*FieldMapping{}
  mappings := &MappingFile{}
  if err := jsoniter.Unmarshal(buf, mappings); err != nil {
    return nil, err
  }
  return mappings, nil
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
  for _, fMap := range maps {
    if strings.ToLower(fMap.SqlType) == strings.ToLower(this.GetSqlType()) {
      deepcopier.Copy(fMap).To(this)
      log.Trace().Str("go type", this.GoType).Str("json type", this.JsonType).Msg("after copying")

      break
    }
  }
}
