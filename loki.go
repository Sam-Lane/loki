package loki

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

//Logger is base interface for different loggers
type Logger interface {
	Set(int)
	WriteFile(string) error
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

//Loki basic logger
type Loki struct {
	loglevel int
	color    bool
	logFile  *os.File
}

const (
	//TRACE level logging
	TRACE int = 1
	//DEBUG level logging
	DEBUG int = 2
	//INFO level logging
	INFO int = 3
	//WARN level logging
	WARN int = 4
	//ERROR level logging
	ERROR int = 5
	//FATAL level logging
	FATAL int = 6
)

//New return a new instance of loki logger.
//By default writes to stdout
func New() *Loki {
	return &Loki{
		color:    true,
		loglevel: INFO,
		logFile:  os.Stdout,
	}
}

//Set the logger level defaults to INFO
func (l *Loki) Set(level int) {
	if level >= 1 && level <= 6 {
		l.loglevel = level
	}
}

//WriteFile tell loki where to write log files too
func (l *Loki) WriteFile(path string) error {
	var err error
	l.logFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf(`failed to open "%s" for writing`, path)
	}
	l.color = false
	return nil
}

//Debug write a log with level DEBUG
func (l *Loki) Debug(mesg string, args ...interface{}) {
	if l.loglevel <= DEBUG {
		l.logStatement("DEBUG", mesg, args...)
	}
}

//Info write a log with level INFO
func (l *Loki) Info(mesg string, args ...interface{}) {
	if l.loglevel <= INFO {
		l.logStatement("INFO", mesg, args...)
	}
}

//Warn write a log with level WARN
func (l *Loki) Warn(mesg string, args ...interface{}) {
	if l.loglevel <= WARN {
		l.logStatement("WARN", mesg, args...)
	}
}

//Error write a log with level ERROR
func (l *Loki) Error(mesg string, args ...interface{}) {
	if l.loglevel <= ERROR {
		l.logStatement("ERROR", mesg, args...)
	}
}

//Fatal with a log with level FATAL
func (l *Loki) Fatal(mesg string, args ...interface{}) {
	if l.loglevel <= FATAL {
		l.logStatement("FATAL", mesg, args...)
	}
}

func (l *Loki) logStatement(level, mesg string, args ...interface{}) {
	level = l.wrapWithColor(level)
	m := fmt.Sprintf(mesg, args...)
	writer := bufio.NewWriter(l.logFile)
	t := time.Now()
	logm := fmt.Sprintf("%s [%s] %s\n", t.Format("2006/01/02 15:01:05"), level, m)
	_, err := writer.WriteString(logm)
	if err != nil {
		fmt.Println("failed to write to buffer")
	}
	writer.Flush()
}

func (l *Loki) wrapWithColor(mesg string) string {
	if isatty.IsTerminal(os.Stdout.Fd()) && l.color {
		return fmt.Sprintf("\033[35m%s\033[0m", mesg)
	}
	return mesg
}
