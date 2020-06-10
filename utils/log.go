package utils

import (
  "github.com/rs/zerolog/log"
)

func IfErrCallback(err error, cb func() interface{}) {
  if err != nil {
    log.Error().Err(err).Msg("")
    cb()
  }
}
