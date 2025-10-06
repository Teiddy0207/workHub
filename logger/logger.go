package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "sync"
	"context"
)

var (
    Log  *zap.Logger
    once sync.Once
)

func InitLogger() {
    once.Do(func() {
        cfg := zap.NewProductionConfig()
        cfg.OutputPaths = []string{"stdout", "logs/app.log"}
        cfg.ErrorOutputPaths = []string{"stderr", "logs/app.log"}
        cfg.EncoderConfig.TimeKey = "time"
        cfg.EncoderConfig.MessageKey = "msg"
        cfg.EncoderConfig.LevelKey = "level"
        cfg.EncoderConfig.CallerKey = "caller"
        cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
        cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
        cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

        logger, err := cfg.Build()
        if err != nil {
            panic(err)
        }
        Log = logger
    })
}

func Info(module, function, message string, fields ...zap.Field) {
    Log.Info(module+":"+function, append([]zap.Field{zap.String("msg", message)}, fields...)...)
}

func Warn(module, function, message string, fields ...zap.Field) {
    Log.Warn(module+":"+function, append([]zap.Field{zap.String("msg", message)}, fields...)...)
}

func Error(module, function, message string, fields ...zap.Field) {
    Log.Error(module+":"+function, append([]zap.Field{zap.String("msg", message)}, fields...)...)
}

func Debug(module, function, message string, fields ...zap.Field) {
    Log.Debug(module+":"+function, append([]zap.Field{zap.String("msg", message)}, fields...)...)
}

func WithTrace(ctx context.Context, module, function string) *zap.Logger {
    var trace string
    if v := ctx.Value("trace_id"); v != nil {
        if s, ok := v.(string); ok {
            trace = s
        }
    }
    if trace == "" {
        trace = "no-trace"
    }
    return Log.With(
        zap.String("module", module),
        zap.String("function", function),
        zap.String("trace_id", trace),
    )
}
