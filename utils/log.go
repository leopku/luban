package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Callback func()

type ConfigOption struct {
	V    bool
	VV   bool
	VVV  bool
	VVVV bool
}

type ConfigOptionFunc func(opts *ConfigOption)

var defaultConfigOption = ConfigOption{
	V:    false,
	VV:   false,
	VVV:  false,
	VVVV: false,
}

func WithV(v bool) ConfigOptionFunc {
	return func(opts *ConfigOption) {
		opts.V = v
	}
}

func WithVV(vv bool) ConfigOptionFunc {
	return func(opts *ConfigOption) {
		opts.VV = vv
	}
}

func WithVVV(vvv bool) ConfigOptionFunc {
	return func(opts *ConfigOption) {
		opts.VVV = vvv
	}
}

func WithVVVV(vvvv bool) ConfigOptionFunc {
	return func(opts *ConfigOption) {
		opts.VVVV = vvvv
	}
}

func InitConfig(opts ...ConfigOptionFunc) {
	options := defaultConfigOption
	for _, opt := range opts {
		opt(&options)
	}

	level := zerolog.ErrorLevel
	if options.V {
		level = zerolog.WarnLevel
	}
	if options.VV {
		level = zerolog.InfoLevel
	}
	if options.VVV {
		level = zerolog.DebugLevel
	}
	if options.VVVV {
		level = zerolog.TraceLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func IfErrCallback(err error, cb func() interface{}) interface{} {
	if err != nil {
		log.Error().Err(err).Msg("")
		return cb()
	}
	return nil
}
