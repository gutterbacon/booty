package zaplog

import (
	"errors"
	"go.amplifyedge.org/booty-v2/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"os"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
)

const (
	DEBUG = zapcore.DebugLevel
	INFO  = zapcore.InfoLevel
	WARN  = zapcore.WarnLevel
)

// ZapLogger implements Logger interface as defined in
// "go.amplifyedge.org/sys-share-v2/sys-core/service/logging"
type ZapLogger struct {
	isDevelopmentMode bool
	logLevel          zapcore.Level
	moduleName        string
	sugarLogger       *zap.SugaredLogger
}

// Zap Logger constructor
func NewZapLogger(level zapcore.Level, moduleName string, isDevelopmentMode bool) *ZapLogger {
	return &ZapLogger{logLevel: level, isDevelopmentMode: isDevelopmentMode, moduleName: moduleName}
}

func (l *ZapLogger) GetLogPath() (string, error) {
	return "", errors.New("")
}

// Init logger
func (l *ZapLogger) InitLogger(extraFields map[string]interface{}) {
	// set log level
	logLevel := l.logLevel
	var logWriter zapcore.WriteSyncer
	var encoderCfg zapcore.EncoderConfig
	logWriter = zapcore.AddSync(os.Stderr)
	encoderCfg = zap.NewDevelopmentEncoderConfig()

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "lvl"
	encoderCfg.CallerKey = "caller"
	encoderCfg.TimeKey = "time"
	encoderCfg.NameKey = "name"
	encoderCfg.MessageKey = "msg"

	if l.isDevelopmentMode {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var zapFields []zap.Field
	zapFields = append(zapFields, zap.Any("app", l.moduleName))
	if extraFields != nil {
		for k, v := range extraFields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel)).With(zapFields)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	_ = l.sugarLogger.Sync()
}

// ZapLogger methods to satisfy Logger interface
func (l *ZapLogger) WithFields(args map[string]interface{}) logging.Logger {
	zl := l
	zl.InitLogger(args)
	return zl
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *ZapLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *ZapLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

func (l *ZapLogger) GetServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	zapOpts := []grpcZap.Option{
		grpcZap.WithLevels(grpcZap.DefaultCodeToLevel),
	}
	return grpcZap.UnaryServerInterceptor(l.sugarLogger.Desugar(), zapOpts...)
}

func (l *ZapLogger) GetServerStreamInterceptor() grpc.StreamServerInterceptor {
	zapOpts := []grpcZap.Option{
		grpcZap.WithLevels(grpcZap.DefaultCodeToLevel),
	}
	return grpcZap.StreamServerInterceptor(l.sugarLogger.Desugar(), zapOpts...)
}
