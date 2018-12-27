package logger

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
)

type LogLevel int

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

// String implements the stringer interface
func (l LogLevel) String() string {
	switch l {
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	}
	return "UNKNOWN"
}

var ErrInvalidLogLevel = errors.New("invalid loglevel")

func StringToLogLevel(s string) (LogLevel, error) {
	switch strings.ToUpper(s) {
	case "FATAL":
		return FATAL, nil
	case "ERROR":
		return ERROR, nil
	case "WARN":
		return WARN, nil
	case "INFO":
		return INFO, nil
	case "DEBUG":
		return DEBUG, nil
	}

	return FATAL, ErrInvalidLogLevel
}

// Logger interface specification
type Logger interface {
	Log(LogLevel, ...interface{})
}

// StdLog turns a standard logger into a Logger interface implementation
type StdLog struct {
	logger      *log.Logger
	Level       LogLevel
	ExitOnFatal bool
}

// Log will output a log line with the caller function, line and the arguments
// Arguments are primarily considered in key/value pairs, i.e. Log("key1", "value1", "key2", "value2")
// If there is no 2nd value argument given, the single argument will be used, i.e. Log("Single argument") or
// Log("key1", "value1", "single argument")
func (l StdLog) Log(lvl LogLevel, args ...interface{}) {
	if lvl > l.Level {
		return
	}

	var buffer bytes.Buffer

	// Get the function name and code line that invoked the log function
	_, fn, line, _ := runtime.Caller(1)

	// Write the name and code line to the log string buffer
	buffer.WriteString(fmt.Sprintf("[%s] %s:%v", lvl, fn, line))

	// Get the amount of arguments to determine the maximum for the for loop
	max := len(args)

	// Loop over all arguments
	i := 0
	for i < max {
		if i < max-1 {
			// We have a next argument, so combine current and next argument as key/value pair
			buffer.WriteString(fmt.Sprintf(" %v: %v", args[i], args[i+1]))
			i += 2
		} else {
			// No next argument, so use single argument
			buffer.WriteString(fmt.Sprintf(" %v", args[i]))
			i++
		}
	}

	// Print the constructed log line
	txt := buffer.String()
	if (lvl <= FATAL) && l.ExitOnFatal {
		l.logger.Fatalf(txt)
	} else {
		l.logger.Printf(txt)
	}
}

func NewStdLog(lvl LogLevel, logger *log.Logger) *StdLog {
	return &StdLog{Level: lvl, logger: logger, ExitOnFatal: true}
}
