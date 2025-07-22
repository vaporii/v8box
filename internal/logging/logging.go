package logging

import (
	"fmt"
)

type LogLevel int

var logLevel LogLevel = LogLevelWarning

const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelVerbose
)

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func Error(format string, a ...any) {
	if logLevel >= LogLevelError {
		fmt.Printf("\u001b[31;1merror: \033[0m"+format+"\n", a...)
	}
}

func Warning(format string, a ...any) {
	if logLevel >= LogLevelWarning {
		fmt.Printf("\u001b[33;1mwarning: \033[0m"+format+"\n", a...)
	}
}

func Info(format string, a ...any) {
	if logLevel >= LogLevelInfo {
		fmt.Printf("\u001b[36;1minfo: \033[0m"+format+"\n", a...)
	}
}

func Verbose(format string, a ...any) {
	if logLevel >= LogLevelVerbose {
		fmt.Printf("\u001b[37mverbose: \033[0m"+format+"\n", a...)
	}
}
