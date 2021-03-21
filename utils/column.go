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
	"github.com/rs/zerolog/log"
	"github.com/segmentio/encoding/json"
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
	if err := json.Unmarshal(buf, mappingFile); err != nil {
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

func (self *ColumnMeta) GetName() string {
	return self.Column.Name()
}

func (self *ColumnMeta) GetGoName() string {
	ret := strcase.ToCamel(self.GetName())
	return ret
}

func (self *ColumnMeta) GetSqlType() string {
	return self.Column.DatabaseTypeName()
}

func (self *ColumnMeta) GetJsonName() string {
	return strcase.ToLowerCamel(self.GetName())
}

func (self *ColumnMeta) ParseAllTypes(maps []*FieldMapping) {
	for i, fMap := range maps {
		found := false
		if strings.ToLower(fMap.SqlType) == strings.ToLower(self.GetSqlType()) {
			deepcopier.Copy(fMap).To(self)
			log.Trace().Str("sql type", self.SqlType).Str("go type", self.GoType).Str("json type", self.JsonType).Msg("after copying")

			found = true
		}

		if found {
			break
		}

		if !found && i == len(maps) {
			self.SqlType = self.GetSqlType()
			self.GoType = "interface{}"
		}
	}
}

// type JenniferFn func(state *jen.Statement) *jen.Statement {}

func (self *ColumnMeta) GenerateGo(group *jen.Group) {
	if strUtil.IsBlank(self.GoType) {
		return
	}
	s := group.Id(self.GetGoName())
	nullable, ok := self.Column.Nullable()
	// jsonTagMap := map[string]string{"json": self.GetJsonName()}
	jsonTag := self.GetJsonName()
	if ok && nullable {
		jsonTag = jsonTag + ",omitempty"
	}

	switch self.GoType {
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

	s.Tag(map[string]string{"json": jsonTag, "db": self.GetName()})

}
