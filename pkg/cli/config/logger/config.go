package logger

import (
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v3"
)

type Config struct {
	logLevel         string
	logFormat        string
	enableStacktrace bool
	enableSource     bool
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Log level (trace, debug, info, warn, error, fatal, panic)",
			Category:    "logger",
			Destination: &x.logLevel,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_LOG_LEVEL")),
			Value:       "info",
		},
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "Log format (json, console)",
			Category:    "logger",
			Destination: &x.logFormat,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_LOG_FORMAT")),
			Value:       "console",
		},
		&cli.BoolFlag{
			Name:        "enable-stacktrace",
			Usage:       "Enable stacktrace in log",
			Category:    "logger",
			Destination: &x.enableStacktrace,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_ENABLE_STACKTRACE")),
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "enable-source",
			Usage:       "Enable source code location in log",
			Category:    "logger",
			Destination: &x.enableSource,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_ENABLE_SOURCE")),
			Value:       false,
		},
	}
}

func (x *Config) LogLevel() string {
	return x.logLevel
}

func (x *Config) LogFormat() string {
	return x.logFormat
}

func (x *Config) EnableStacktrace() bool {
	return x.enableStacktrace
}

func (x *Config) Build() (*slog.Logger, error) {
	output := os.Stdout

	logLevelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	level, ok := logLevelMap[x.logLevel]
	if !ok {
		return nil, goerr.New("invalid log level").With("level", x.logLevel)
	}

	var handler slog.Handler
	switch x.logFormat {
	case "json":
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: x.enableSource,
			Level:     level,
		})

	case "console":
		options := []clog.Option{
			clog.WithWriter(output),
			clog.WithLevel(level),
			// clog.WithReplaceAttr(filter),

			// clog.WithTimeFmt("2006-01-02 15:04:05"),
			clog.WithColorMap(&clog.ColorMap{
				Level: map[slog.Level]*color.Color{
					slog.LevelDebug: color.New(color.FgGreen, color.Bold),
					slog.LevelInfo:  color.New(color.FgCyan, color.Bold),
					slog.LevelWarn:  color.New(color.FgYellow, color.Bold),
					slog.LevelError: color.New(color.FgRed, color.Bold),
				},
				LevelDefault: color.New(color.FgBlue, color.Bold),
				Time:         color.New(color.FgWhite),
				Message:      color.New(color.FgHiWhite),
				AttrKey:      color.New(color.FgHiCyan),
				AttrValue:    color.New(color.FgHiWhite),
			}),
		}

		if x.enableSource {
			options = append(options, clog.WithSource(true))
		}

		if x.enableStacktrace {
			options = append(options, clog.WithAttrHook(clog.GoerrHook))
		} else {
			options = append(options, clog.WithAttrHook(defaultErrHook))
		}

		handler = clog.New(options...)

	default:
		return nil, goerr.New("invalid log format").With("format", x.logFormat)
	}

	return slog.New(handler), nil
}

func defaultErrHook(_ []string, attr slog.Attr) *clog.HandleAttr {
	goErr, ok := attr.Value.Any().(*goerr.Error)
	if !ok {
		return nil
	}

	var attrs []any
	for k, v := range goErr.Values() {
		attrs = append(attrs, slog.Any(k, v))
	}
	newAttr := slog.Group(attr.Key,
		slog.String("msg", goErr.Error()),
		slog.Group("attrs", attrs...),
	)

	return &clog.HandleAttr{
		NewAttr: &newAttr,
	}
}
