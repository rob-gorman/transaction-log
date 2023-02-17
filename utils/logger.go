package utils

import (
	"log"
	"os"
)

type AuditLogLogger struct {
	info *log.Logger
	err  *log.Logger
}

func NewLogger() *AuditLogLogger {
	return &AuditLogLogger{
		info: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		err:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *AuditLogLogger) Info(s string, vals ...interface{}) {
	l.info.Printf(s, vals...)
}

func (l *AuditLogLogger) Err(s string, vals ...interface{}) {
	l.err.Printf(s, vals...)
}
