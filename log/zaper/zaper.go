// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz
// Create Time 2020/12/7

package zaper

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger strcut
type logger struct {
	errorPath string
	infoPath  string
	maxAge    time.Duration
	rotaTime  time.Duration
	atom      zap.AtomicLevel
	sugar     *zap.SugaredLogger
}

func (l *logger) getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename,
		//rotatelogs.WithLinkName(filename),
		//rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithMaxAge(l.maxAge),
		//rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationTime(l.rotaTime),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func New(errorPath, infoPath string, maxAge, rotaTime time.Duration) (*logger, error) {

	if errorPath == "" {
		errorPath = "./logs/app_error.log"
	}
	errorPath = strings.Replace(errorPath, ".log", "", -1) + "-%Y%m%d%H.log"

	if infoPath == "" {
		infoPath = "./logs/app_info.log"
	}
	infoPath = strings.Replace(infoPath, ".log", "", -1) + "-%Y%m%d%H.log"

	log := &logger{}
	log.errorPath = errorPath
	log.infoPath = infoPath
	log.maxAge = maxAge
	log.rotaTime = rotaTime
	log.atom = zap.NewAtomicLevel()

	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	//encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	// debug < info < warn < error < panic < fatal
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := log.getWriter(infoPath)
	errorWriter := log.getWriter(errorPath)
	// 创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	zaplog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	log.sugar = zaplog.Sugar()

	return log, nil
}

func (l *logger) SetLevel(level string) {

	var zaplevel zapcore.Level

	switch level {
	case "debug":
		zaplevel = zap.DebugLevel
	case "info":
		zaplevel = zap.InfoLevel
	case "warn":
		zaplevel = zap.WarnLevel
	case "error":
		zaplevel = zap.ErrorLevel
	default:
		zaplevel = zap.InfoLevel
	}

	l.atom.SetLevel(zaplevel)
}
func (l *logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}
func (l *logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}
func (l *logger) TDebug(tag, message string, args ...interface{}) {
	l.sugar.Debugf(tag+"\t"+message, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}
func (l *logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}
func (l *logger) TInfo(tag, message string, args ...interface{}) {
	l.sugar.Infof(tag+"\t"+message, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}
func (l *logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}
func (l *logger) TWarn(tag, message string, args ...interface{}) {
	l.sugar.Warnf(tag+"\t"+message, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}
func (l *logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}
func (l *logger) TError(tag, message string, args ...interface{}) {
	l.sugar.Errorf(tag+"\t"+message, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}
func (l *logger) DPanicf(template string, args ...interface{}) {
	l.sugar.DPanicf(template, args...)
}
func (l *logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}
func (l *logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}
func (l *logger) TPanic(tag, message string, args ...interface{}) {
	l.sugar.DPanicf(tag+"\t"+message, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}
func (l *logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}
func (l *logger) TFatal(tag, message string, args ...interface{}) {
	l.sugar.Fatalf(tag+"\t"+message, args...)
}

func (l *logger) Close() error {
	return l.sugar.Sync()
}
