package utils

import (
  "io/ioutil"

  "github.com/json-iterator/go"
  // "github.com/rs/zerolog/log"
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
