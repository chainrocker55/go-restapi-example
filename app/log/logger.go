package log

import (
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/consts"

	"log/slog"
	"os"
	"strings"

	"github.com/modern-go/gls"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logEnvKey             string = "env"
	logTypeKey            string = "logtype"
	logApplicationNameKey string = "application.name"
	logTypeApplication    string = "application"
	logMessageKey         string = "message"
	logTimestampKey       string = "timestamp"
	logTimeFormat                = "2006-01-02 15:04:05.000"
)

var CloseLogger func() = func() {
	// Set empty function as default close logger function
}

var logger *zap.Logger = zap.L()

func GetLogger() *zap.Logger {
	return logger
}

func Init() *zap.Logger {

	closeLogstashConnection := func() error { return nil }

	var loggingFeed = configs.Conf.Logger.Feed
	var logLevel = configs.Conf.Logger.Level
	var loggingUrl = configs.Conf.Logger.Url

	zapLogLevel := getLogLevel(logLevel)
	zapDevConfig := zap.NewDevelopmentEncoderConfig()

	zapDevConfig.EncodeTime = zapcore.TimeEncoderOfLayout(logTimeFormat)
	consoleEncoder := zapcore.NewConsoleEncoder(zapDevConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), zapLogLevel)

	var logCores = []zapcore.Core{consoleCore}
	if loggingFeed {

		logstashConn, err := NewLogstashConnection(loggingUrl)
		if err != nil {
			slog.Warn("connecting to Logstash failed. Logger will try to connect next time logger is used", slog.Any("error", err))
		}

		closeLogstashConnection = logstashConn.Close
		zapProductConfig := zap.NewProductionEncoderConfig()
		zapProductConfig.MessageKey = logMessageKey
		zapProductConfig.TimeKey = logTimestampKey
		tcpEncoder := zapcore.NewJSONEncoder(zapProductConfig)
		tcpCore := zapcore.NewCore(tcpEncoder, zapcore.AddSync(logstashConn), zapLogLevel)

		logCores = append(logCores, tcpCore)
	}

	teeCore := zapcore.NewTee(logCores...)

	logger = zap.New(teeCore).With(
		zap.String(logApplicationNameKey, configs.Conf.App.Name),
		zap.String(logEnvKey, configs.Conf.App.Env),
		zap.String(logTypeKey, logTypeApplication),
	)

	slog.Info("Init log success!")

	undoReplacedGlobalLog := zap.ReplaceGlobals(logger)
	slog.Info("Init log replace globals success!")

	CloseLogger = func() {
		logger.Sync()
		undoReplacedGlobalLog()
		closeLogstashConnection()
	}

	return logger
}

func getLogLevel(logLevel string) zapcore.Level {

	switch strings.ToLower(logLevel) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func wrapFields() *zap.Logger {
	xCorrelationId := getContextValue(consts.ContextCorrelationId)
	xRequestId := getContextValue(consts.ContextRequestId)

	return logger.With(
		zap.String(consts.ContextCorrelationId, xCorrelationId),
		zap.String(consts.ContextRequestId, xRequestId),
	)
}

func Debug(msg string) {
	wrapFields().Debug(msg)
}

func Debugf(template string, args ...interface{}) {
	wrapFields().Sugar().Debugf(template, args...)
}

func Info(msg string) {
	wrapFields().Info(msg)
}

func Infof(template string, args ...interface{}) {
	wrapFields().Sugar().Infof(template, args...)
}

func Warn(msg string) {
	wrapFields().Warn(msg)
}

func Warnf(template string, args ...interface{}) {
	wrapFields().Sugar().Warnf(template, args...)
}

func Error(args ...interface{}) {
	wrapFields().Sugar().Error(args...)
}

func Errorf(template string, args ...interface{}) {
	wrapFields().Sugar().Errorf(template, args...)
}

func Panic(args ...interface{}) {
	wrapFields().Sugar().Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	wrapFields().Sugar().Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	wrapFields().Sugar().Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	wrapFields().Sugar().Fatalf(template, args...)
}

func getContextValue(key string) string {
	value := gls.Get(key)
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return ""
}
