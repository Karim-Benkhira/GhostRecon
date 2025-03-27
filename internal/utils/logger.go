package utils

import (
	"fmt"
	"time"
)


type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)


type Logger struct {
	Level LogLevel
}


func NewLogger(level LogLevel) *Logger {
	return &Logger{
		Level: level,
	}
}


func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	if level >= l.Level {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		levelStr := []string{"DEBUG", "INFO", "WARN", "ERROR"}[level]
		message := fmt.Sprintf(format, args...)
		fmt.Printf("[%s] %s: %s\n", timestamp, levelStr, message)
	}
}
