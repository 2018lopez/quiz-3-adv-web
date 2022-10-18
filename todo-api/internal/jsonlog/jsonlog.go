//Filename: internal/jsonlog/jsonlog.go

package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// difference severity levels of logging entries
type Level int8

// Levels start at zero
const (
	LevelInfo  Level = iota // value is 0
	LevelError              // value is 1
	LevelFatal              // value is 2
	LevelOff                // value is 3
)

// Severity level as a human readeable friendly format
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "Info"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "Fatal"
	default:
		return ""
	}
}

//define custome logger

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// New() function create a new instance of logger

func New(out io.Writer, minLevel Level) *Logger {

	return &Logger{
		out:      out,
		minLevel: minLevel,
	}

}

// hepler methods
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	//Ensure the severity level is minimum
	if level < l.minLevel {
		return 0, nil
	}
	//creat struct for data

	data := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	//include the stack trace
	if level >= LevelError {
		data.Trace = string(debug.Stack())
	}

	var entry []byte
	entry, err := json.Marshal(data)
	if err != nil {
		entry = []byte(LevelError.String() + " : unable to marshal log message: " + err.Error())
	}
	//write the log entry
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(entry, '\n'))
}

// implement the io.writer interface
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
