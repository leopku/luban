package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	// strUtil "github.com/agrison/go-commons-lang/stringUtils"
	// "github.com/dave/jennifer/jen"
	"github.com/fatih/structs"
	"github.com/rs/zerolog/log"
	ramda "github.com/vijaykanthm28/go-ramda"
)

func NewConfigFromReader(in io.Reader, cType string) *viper.Viper {
	cfg := viper.New()
	cfg.SetConfigType(cType)
	cfg.ReadConfig(in)
	log.Trace().Interface("cfg", cfg).Msg("")
	return cfg
}

func GetModelPath(outpath string) string {
	// ret := vip.GetString("generation.model.output")
	ret := outpath
	base := path.Base(ret)
	if base == "/" || base == "." {
		ret = "./models"
	}
	return ret
}

func GetModelPackageName(outpath string) string {
	return path.Base(GetModelPath(outpath))
}

func CreateDirectory(dirName string) error {
	src, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return err
		}
		return nil
	}

	if src.Mode().IsRegular() {
		return errors.New(fmt.Sprintf("%s already exist as a file", dirName))
	}
	return nil
}

func Keys(t interface{}, exclude []string) []string {
	s := structs.New(t)
	var keys []string
	for _, field := range s.Fields() {
		key := strings.Split(field.Tag("db"), ",")[0]
		if len(key) > 0 && key != "-" && ramda.Contains(exclude, key) {
			keys = append(keys, key)
		}
	}
	return keys
}
