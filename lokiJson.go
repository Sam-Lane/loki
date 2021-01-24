package loki

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

//LokiJson logs in json format
type LokiJson struct {
	loglevel int
	logFile  *os.File
}

//JsonLog struct defines the structure of a log message
type JsonLog struct {
	Timestamp time.Time          `json:"timestamp"`
	Message   string             `json:"message"`
	Error     error              `json:"error,omitempty"`
	Level     string             `json:"level"`
	Caller    Caller             `json:"caller"`
	Context   []*json.RawMessage `json:"context,omitempty"`
}

//Caller defines json structure to store data on where the log message orginated
type Caller struct {
	Function string `json:"function"`
	Line     int    `json:"line"`
	File     string `json:"file"`
}

//NewJsonLogger return a new loki logger that logs in json format.
//By default writes to stdout
func NewJsonLogger() *LokiJson {
	return &LokiJson{
		loglevel: INFO,
		logFile:  os.Stdout,
	}
}

//Set the logging level
func (l *LokiJson) Set(level int) {
	if level >= 1 && level <= 6 {
		l.loglevel = level
	} else {
		l.loglevel = INFO
	}
}

func (l *LokiJson) WriteFile(path string) error {
	var err error
	l.logFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf(`failed to open "%s" for writing`, path)
	}
	return nil
}

//Debug write a log with level DEBUG
func (l *LokiJson) Debug(mesg string, args ...interface{}) {
	if l.loglevel <= DEBUG {
		l.logMessage(mesg, "DEBUG", args...)
	}
}

//Info write a log with level INFO
func (l *LokiJson) Info(mesg string, args ...interface{}) {
	if l.loglevel <= INFO {
		l.logMessage(mesg, "INFO", args...)
	}
}

//Warn write a log with level WARN
func (l *LokiJson) Warn(mesg string, args ...interface{}) {
	if l.loglevel <= WARN {
		l.logMessage(mesg, "WARN", args...)
	}
}

//Error write a log with level ERROR
func (l *LokiJson) Error(mesg string, args ...interface{}) {
	if l.loglevel <= ERROR {
		l.logMessage(mesg, "ERROR", args...)
	}
}

//Fatal write a log with level FATAL
func (l *LokiJson) Fatal(mesg string, args ...interface{}) {
	if l.loglevel <= FATAL {
		l.logMessage(mesg, "FATAL", args...)
	}
}

func (l *LokiJson) logMessage(mesg, level string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()
	j := JsonLog{
		Timestamp: time.Now(),
		Message:   mesg,
		Level:     level,
		Caller: Caller{
			Function: function,
			Line:     line,
			File:     file,
		},
	}

	writer := bufio.NewWriter(l.logFile)
	_, err := writer.WriteString(l.structToJson(&j))
	if err != nil {
		fmt.Println("failed to write to buffer")
	}
	writer.Flush()
}

func (l *LokiJson) structToJson(s *JsonLog) string {
	b, _ := json.Marshal(s)
	return string(b) + "\n"
}
