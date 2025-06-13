package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

func NewLogger(filePath string) (*Logger, error) {
	// file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		//   file:        file,
	}, nil
}

func (l *Logger) Info(message string, interfaceArgs ...interface{}) {
	if len(interfaceArgs) > 0 {
		l.infoLogger.Printf(message, interfaceArgs...)
	} else {
		l.infoLogger.Print(message)
	}
}

func (l *Logger) Error(message string, interfaceArgs ...interface{}) {
	if len(interfaceArgs) > 0 {
		l.errorLogger.Printf(message, interfaceArgs...)
	} else {
		l.errorLogger.Print(message)
	}
}

func (l *Logger) Close() {
	l.file.Close()
}
