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

func (l *Logger) Info(message string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLogger.Printf(message, args...)
	} else {
		l.infoLogger.Printf("%s", message)
	}
}

func (l *Logger) Error(message string, err error) {
	l.errorLogger.Printf("%s: %v", message, err)
}

func (l *Logger) Close() {
	l.file.Close()
}
