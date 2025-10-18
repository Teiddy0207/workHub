package logger

import (
	"sync"
	"time"

	"workHub/internal/config"
	"workHub/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const RequestContextKey = "request_context"

// RequestContext holds request-specific logging information
type RequestContext struct {
	Timestamp time.Time
	Method    string
	Path      string
	RequestID string
	UserID    string
	ClientIP  string // Thêm trường ClientIP
	Error     error
}

// Logger interface for common logging operations
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

// Global variables
var (
	logger                *zap.Logger
	currentRequestContext *RequestContext
	contextMutex          sync.RWMutex
	L                     Logger
)

// Log levels mapping
var logLevels = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

// defaultLogger implements Logger interface
type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, fields ...zap.Field) {
	if ctxFields := getContextFields(); ctxFields != nil {
		fields = append(fields, ctxFields...)
	}
	logger.Debug(msg, fields...)
}

func (l *defaultLogger) Info(msg string, fields ...zap.Field) {
	if ctxFields := getContextFields(); ctxFields != nil {
		fields = append(fields, ctxFields...)
	}
	logger.Info(msg, fields...)
}

func (l *defaultLogger) Warn(msg string, fields ...zap.Field) {
	if ctxFields := getContextFields(); ctxFields != nil {
		fields = append(fields, ctxFields...)
	}
	logger.Warn(msg, fields...)
}

func (l *defaultLogger) Error(msg string, fields ...zap.Field) {
	if ctxFields := getContextFields(); ctxFields != nil {
		fields = append(fields, ctxFields...)
	}
	logger.Error(msg, fields...)
}

func (l *defaultLogger) Fatal(msg string, fields ...zap.Field) {
	if ctxFields := getContextFields(); ctxFields != nil {
		fields = append(fields, ctxFields...)
	}
	logger.Fatal(msg, fields...)
}

// Context management functions
func SetRequestContext(reqCtx *RequestContext) {
	contextMutex.Lock()
	defer contextMutex.Unlock()
	currentRequestContext = reqCtx
}

func GetCurrentRequestContext() *RequestContext {
	contextMutex.RLock()
	defer contextMutex.RUnlock()
	return currentRequestContext
}

func ClearRequestContext() {
	contextMutex.Lock()
	defer contextMutex.Unlock()
	currentRequestContext = nil
}

// Get context fields for logging
func getContextFields() []zap.Field {
	reqCtx := GetCurrentRequestContext()
	if reqCtx == nil {
		return nil
	}

	fields := make([]zap.Field, 0, 5) // Tăng capacity lên 5

	if reqCtx.Method != "" {
		fields = append(fields, zap.String("method", reqCtx.Method))
	}
	if reqCtx.Path != "" {
		fields = append(fields, zap.String("path", reqCtx.Path))
	}
	if reqCtx.RequestID != "" {
		fields = append(fields, zap.String("request_id", reqCtx.RequestID))
	}
	if reqCtx.UserID != "" {
		fields = append(fields, zap.String("user_id", reqCtx.UserID))
	}
	if reqCtx.ClientIP != "" {
		fields = append(fields, zap.String("client_ip", reqCtx.ClientIP))
	}

	return fields
}

// Package-level logging functions
func Debug(msg string, fields ...zap.Field) {
	L.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	L.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	L.Warn(msg, fields...)
}

// Package-level logging functions đơn giản hóa
func Error(msg string, err error) {
	if reqCtx := GetCurrentRequestContext(); reqCtx != nil {
		reqCtx.Error = err // Chỉ set để tracking
	}
	L.Error(msg, zap.Error(err)) // Chỉ log error ở đây
}

func Fatal(msg string, err error) {
	if reqCtx := GetCurrentRequestContext(); reqCtx != nil {
		reqCtx.Error = err // Chỉ set để tracking
	}
	L.Fatal(msg, zap.Error(err)) // Chỉ log error ở đây
}

// Write syncer configuration
func getWriteSyncer() zapcore.WriteSyncer {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/vnpt-be.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	return zapcore.NewMultiWriteSyncer(fileWriter)
}

// Initialize creates and configures the logger instance
func Initialize(cfg config.Logger, nodeEnv constant.NodeEnv) {
	// Determine log level
	logLevel := zapcore.InfoLevel
	if level, exists := logLevels[cfg.Level]; exists {
		logLevel = level
	}

	// Configure encoder based on environment
	var encoderCfg zapcore.EncoderConfig
	if nodeEnv == constant.NODE_ENV_DEVELOPMENT {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	// Customize encoder config
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "msg"
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	// Create logger
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		getWriteSyncer(),
		zap.NewAtomicLevelAt(logLevel),
	)

	logger = zap.New(core)
	L = &defaultLogger{}
}
