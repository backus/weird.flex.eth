package main

import (
	"fmt"
	"strings"
)

type LogLevel byte

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type Logger struct {
	level LogLevel
}

func NewLogger() Logger {
	return Logger{LogLevelDebug}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l Logger) Debug(message string, args ...interface{}) {
	handleLog(l, LogLevelDebug, message, args...)
}

func (l Logger) Info(message string, args ...interface{}) {
	handleLog(l, LogLevelInfo, message, args...)
}

func (l Logger) Warn(message string, args ...interface{}) {
	handleLog(l, LogLevelWarn, message, args...)
}

func (l Logger) Error(message string, args ...interface{}) {
	handleLog(l, LogLevelError, message, args...)
}

func (l Logger) Fatal(message string, args ...interface{}) {
	handleLog(l, LogLevelFatal, message, args...)
}

func handleLog(logger Logger, level LogLevel, message string, args ...interface{}) {
	if logger.level > level {
		return
	}

	msg := message

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	fmt.Printf(msg, args...)
}
