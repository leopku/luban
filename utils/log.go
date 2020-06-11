package utils

import (
  "github.com/rs/zerolog/log"
)

type Callback func()

func IfErrCallback(err error, cb func() interface{}) interface{} {
  if err != nil {
    log.Error().Err(err).Msg("")
    return cb()
  }
  return nil
}
