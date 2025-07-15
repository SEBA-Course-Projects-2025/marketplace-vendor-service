package logs

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type LokiGormLogger struct {
	labels map[string]string
	level  logger.LogLevel
}

func NewLokiGormLogger(newLabels map[string]string, newLevel logger.LogLevel) *LokiGormLogger {
	return &LokiGormLogger{labels: newLabels, level: newLevel}
}

func (lgl *LokiGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	lgl.level = level
	return lgl
}

func (lgl *LokiGormLogger) Info(ctx context.Context, message string, data ...interface{}) {
	if lgl.level >= logger.Info {
		select {
		case logsQueue <- logEntry{fmt.Sprintf("INFO: "+message, data...), lgl.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("INFO: "+message, data...))
		}
	}
}

func (lgl *LokiGormLogger) Warn(ctx context.Context, message string, data ...interface{}) {
	if lgl.level >= logger.Warn {
		select {
		case logsQueue <- logEntry{fmt.Sprintf("WARN: "+message, data...), lgl.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("WARN: "+message, data...))
		}
	}
}

func (lgl *LokiGormLogger) Error(ctx context.Context, message string, data ...interface{}) {
	if lgl.level >= logger.Error {
		select {
		case logsQueue <- logEntry{fmt.Sprintf("ERROR: "+message, data...), lgl.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+message, data...))
		}
	}
}

func (lgl *LokiGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	elapsedTime := time.Since(begin)

	sql, rows := fc()

	sqlMessage := fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsedTime.Microseconds())/1000.0, rows, sql)

	switch {
	case err != nil && lgl.level >= logger.Error:
		select {
		case logsQueue <- logEntry{"Error executing SQL query:" + sqlMessage, lgl.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+sqlMessage, lgl.labels))
		}
	case lgl.level == logger.Info:
		select {
		case logsQueue <- logEntry{"SQL:" + sqlMessage, lgl.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+sqlMessage, lgl.labels))
		}
	}
}
