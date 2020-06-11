package utils

import (
  "errors"
  "fmt"
  "os"
  "path"

  "github.com/spf13/viper"
  // strUtil "github.com/agrison/go-commons-lang/stringUtils"
)

func GetModelPath() string {
  ret := viper.GetString("generation.model.output")
  base := path.Base(ret)
  if base == "/" || base == "." {
    ret = "./models"
  }
  return ret
}

func GetModelPackageName() string {
  return path.Base(GetModelPath())
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
