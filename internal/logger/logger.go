package logger

import (
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

type LoggerConfig struct {
	Level       string
	Environment string
}

type Field struct {
	Key   string
	Value any
}

type Logger struct {
	log    *zap.Logger
	fields []Field
}

func (c *LoggerConfig) New() (*Logger, error) {
	settings := zap.NewProductionConfig()
	settings.EncoderConfig = ecszap.ECSCompatibleEncoderConfig(settings.EncoderConfig)

	if c.Environment == "local" {
		settings = prettyconsole.NewConfig()
	}

	settings.OutputPaths = []string{"stdout"}
	err := settings.Level.UnmarshalText([]byte(c.Level))
	if err != nil {
		return nil, err
	}

	zapLogger, err := settings.Build(ecszap.WrapCoreOption(), zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{
		log:    zapLogger,
		fields: make([]Field, 0),
	}, nil
}

func (l Logger) Fields(fields ...Field) *Logger {
	l.fields = append(l.fields, fields...)
	return &l
}

func (l Logger) Skip(skip int) *Logger {
	l.log.WithOptions(zap.AddCallerSkip(1))
	return &l
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.fields = append(l.fields, fields...)
	l.log.Debug(msg, parseFields(fields...)...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.fields = append(l.fields, fields...)
	l.log.Info(msg, parseFields(fields...)...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.fields = append(l.fields, fields...)
	l.log.Warn(msg, parseFields(fields...)...)
}

func (l *Logger) Error(err error, msg string, fields ...Field) {
	l.fields = append(l.fields, fields...)
	zapFields := parseFields(l.fields...)
	zapFields = append(zapFields, zap.Error(err))
	l.log.Error(msg, zapFields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.fields = append(l.fields, fields...)
	l.log.Panic(msg, parseFields(fields...)...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.fields = append(l.fields, fields...)
	l.log.Fatal(msg, parseFields(fields...)...)
}

func parseFields(fields ...Field) []zap.Field {
	zapFields := []zap.Field{}
	for _, field := range fields {
		switch field.Value.(type) {
		case int:
			val := field.Value.(int)
			zapFields = append(zapFields, zap.Int(field.Key, val))
		case bool:
			val := field.Value.(bool)
			zapFields = append(zapFields, zap.Bool(field.Key, val))
		case float32:
			val := field.Value.(float32)
			zapFields = append(zapFields, zap.Float32(field.Key, val))
		case float64:
			val := field.Value.(float64)
			zapFields = append(zapFields, zap.Float64(field.Key, val))
		case string:
			val := field.Value.(string)
			zapFields = append(zapFields, zap.String(field.Key, val))
		case error:
			val := field.Value.(error)
			zapFields = append(zapFields, zap.Error(val))
		default:
			zapFields = append(zapFields, zap.Any(field.Key, field.Value))
		}
	}

	return zapFields
}
