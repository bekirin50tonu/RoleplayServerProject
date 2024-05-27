package logger

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warning
	Error
)

type Log struct {
	Time    time.Time
	Level   LogLevel
	Service string
	Message string
}

type LogManager struct {
	mu      sync.Mutex
	logFile *os.File
	logger  *log.Logger
}

func NewLogManager(logFilePath string, output io.Writer) (*LogManager, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &LogManager{
		logFile: logFile,
		logger:  log.New(io.MultiWriter(logFile, output), "", 0),
	}, nil
}

func (lm *LogManager) Log(level LogLevel, service, message string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.logger.Printf("[%s] - [%s] - [%s] - %s\n", time.Now().Format("02.01.2006 - 15:04"), service, levelToString(level), message)
}

func (lm *LogManager) Close() {
	lm.logFile.Close()
}

func levelToString(level LogLevel) string {
	switch level {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type DebugLogger struct {
	LogManager
	Service string
}

func NewDebugLogger(logFilePath string, output io.Writer, service string) (*DebugLogger, error) {
	logManager, err := NewLogManager(logFilePath, output)
	if err != nil {
		return nil, err
	}

	return &DebugLogger{
		LogManager: *logManager,
		Service:    service,
	}, nil
}

func (dl *DebugLogger) Debug(message string) {
	dl.Log(Debug, dl.Service, message)
}
