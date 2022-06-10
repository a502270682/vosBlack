package log

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Entry = logrus.Entry
type Ext1FieldLogger = logrus.Ext1FieldLogger
type FieldLogger = logrus.FieldLogger
type Fields = logrus.Fields
type Formatter = logrus.Formatter
type Hook = logrus.Hook
type Logger = MyLogger
type Level = logrus.Level
type LevelHooks = logrus.LevelHooks
type MutexWrap = logrus.MutexWrap

const PanicLevel = logrus.PanicLevel
const FatalLevel = logrus.FatalLevel
const ErrorLevel = logrus.ErrorLevel
const WarnLevel = logrus.WarnLevel
const InfoLevel = logrus.InfoLevel
const DebugLevel = logrus.DebugLevel
const TraceLevel = logrus.TraceLevel

var logger = setLogger()

func setLogger() Logger {
	formatter := &logrus.JSONFormatter{}
	nl := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	sl := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &CtxLogger{&nl, &sl}
}

type Options struct {
	Level     string `mapstructure:"level" json:"level" toml:"level"`
	File      string `mapstructure:"file" json:"file" toml:"file"`
	ErrFile   string `mapstructure:"err_file" json:"err_file" toml:"err_file"`
	CrashFile string `mapstructure:"crash_file" json:"crash_file" toml:"crash_file"`
	AppName   string `mapstructure:"app_name" json:"app_name" toml:"app_name"`
	Format    string `mapstructure:"format" json:"format" toml:"format"`
}

func NewLoggerWithOptions(options Options) (l Logger, err error) {
	l = setLogger()
	if err = initLoggerWithOptions(l, options); err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}
	return l, nil
}

func initLoggerWithOptions(l Logger, options Options) (err error) {
	if options.Level != "" { // 如果配置里指定了日志等级，则解析并设置，否则默认等级是info。
		level, err := ParseLevel(options.Level)
		if err != nil {
			return errors.Wrapf(err, "failed to parse level(%s)", options.Level)
		}
		l.SetLevel(level)
	}
	if options.File != "" { // 如果配置里指定了日志文件，则解析并设置，否则默认写到stderr。
		err = handleFileOutput(l, options.File) // 设置output、压测标志
		if err != nil {
			return errors.Wrapf(err, "failed to set logger.Output and set flow_control")
		}
	}

	return
}

func handleFileOutput(l Logger, fileName string) error {
	writer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open file(%s)", fileName)
	}
	l.SetOutput(writer) // 设置正常日志
	return nil
}

func GetLogger() Logger {
	return logger
}

func SetOutput(out io.Writer) {
	logger.SetOutput(out)
}

func SetFormatter(formatter Formatter) {
	logger.SetFormatter(formatter)
}

func SetReportCaller(include bool) {
	logger.SetReportCaller(include)
}

func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func SetLevelWithShadow(level logrus.Level) {
	logger.SetLevel(level)
}

func AddHook(hook logrus.Hook) {
	logger.AddHook(hook)
}

func ParseLevel(level string) (Level, error) {
	return logrus.ParseLevel(level)
}

func NewLogrusEntry(l *logrus.Logger) *Entry {
	return logrus.NewEntry(l)
}

func WithField(key string, value interface{}) LoggerEntry {
	return logger.WithField(key, value)
}

func WithFields(fields Fields) LoggerEntry {
	return logger.WithFields(fields)
}

func WithTime(t time.Time) LoggerEntry {
	return logger.WithTime(t)
}

func WithObject(obj interface{}) LoggerEntry {
	return logger.WithObject(obj)
}

func Trace(ctx context.Context, args ...interface{}) {
	logger.Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	logger.Debug(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	logger.Print(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	logger.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	logger.Warn(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	logger.Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	logger.Error(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	logger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	logger.Fatal(ctx, args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	logger.Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(ctx, format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.Warnf(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	logger.Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.Errorf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	logger.Panicf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logger.Fatalf(ctx, format, args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	logger.Traceln(ctx, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	logger.Debugln(ctx, args...)
}

func Println(ctx context.Context, args ...interface{}) {
	logger.Println(ctx, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	logger.Infoln(ctx, args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	logger.Warnln(ctx, args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	logger.Warningln(ctx, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	logger.Errorln(ctx, args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	logger.Panicln(ctx, args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	logger.Fatalln(ctx, args...)
}
