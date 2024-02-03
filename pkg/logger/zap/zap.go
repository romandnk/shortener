package zaplogger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Test             bool     `yaml:"test"`
	Level            string   `yaml:"level"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
}

type Logger struct {
	logger *zap.Logger
}

func NewLogger(cfg Config) (*Logger, error) {
	var (
		err error
		c   *zap.Config
	)

	if cfg.Test {
		c, err = devCfg(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		c = prodCfg(cfg)
		if err != nil {
			return nil, err
		}
	}

	logger, err := c.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger}, nil
}

func devCfg(cfg Config) (*zap.Config, error) {
	lvl, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	return &zap.Config{
		Level:            lvl,
		Encoding:         "console",
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "lvl",
			TimeKey:    "ts",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.DateTime))
			},
			EncodeLevel: func(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(lvl.String())
			},
			ConsoleSeparator: " | ",
		},
	}, nil
}

func prodCfg(cfg Config) *zap.Config {
	return &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "lvl",
			TimeKey:    "ts",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.DateTime))
			},
			EncodeLevel: func(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(lvl.String())
			},
		},
	}
}

func (l *Logger) Info(msg string, args ...any) {
	zapFields := make([]zap.Field, 0, len(args))

	for _, field := range args {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.logger.Info(msg, zapFields...)
}

func (l *Logger) Error(msg string, args ...any) {
	zapFields := make([]zap.Field, 0, len(args))

	for _, field := range args {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.logger.Error(msg, zapFields...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	zapFields := make([]zap.Field, 0, len(args))

	for _, field := range args {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.logger.Fatal(msg, zapFields...)
}
