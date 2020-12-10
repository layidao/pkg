// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz@stc.plus
// Create Time 2020/12/10

package log

// Logger interface
type Provider interface {
	SetLevel(level string)

	Info(args ...interface{})
	Infof(template string, args ...interface{})
	TInfo(tag, template string, args ...interface{})

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	TWarn(tag, template string, args ...interface{})

	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	TDebug(tag, template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	TError(tag, template string, args ...interface{})

	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	TPanic(tag, template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})

	Close() error
}

var Log Provider
