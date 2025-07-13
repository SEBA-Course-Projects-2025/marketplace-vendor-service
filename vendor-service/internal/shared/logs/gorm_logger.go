package logs

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type logEntry struct {
	message string
	labels  map[string]string
}

var logsQueue = make(chan logEntry, 1000)

func init() {

	workers := 8

	for i := 0; i < workers; i++ {
		go func() {
			for entry := range logsQueue {
				if err := SendLogsToLoki(entry.message, entry.labels); err != nil {
					log.Printf("Error sending logs to grafana: %v", err)
				}
			}
		}()
	}

}

type LokiLogger struct {
	labels map[string]string
	level  logger.LogLevel
}

func NewLokiLogger(newLabels map[string]string, newLevel logger.LogLevel) *LokiLogger {
	return &LokiLogger{labels: newLabels, level: newLevel}
}

func (ll *LokiLogger) LogMode(level logger.LogLevel) logger.Interface {
	ll.level = level
	return ll
}

func (ll *LokiLogger) Info(ctx context.Context, message string, data ...interface{}) {
	if ll.level >= logger.Info {
		select {
		case logsQueue <- logEntry{fmt.Sprintf("INFO: "+message, data...), ll.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("INFO: "+message, data...))
		}
	}
}

func (ll *LokiLogger) Warn(ctx context.Context, message string, data ...interface{}) {
	if ll.level >= logger.Warn {
		select {
		case logsQueue <- logEntry{fmt.Sprintf("WARN: "+message, data...), ll.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("WARN: "+message, data...))
		}
	}
}

func (ll *LokiLogger) Error(ctx context.Context, message string, data ...interface{}) {
	if ll.level >= logger.Error {
		select {
		case logsQueue <- logEntry{fmt.Sprintf("ERROR: "+message, data...), ll.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+message, data...))
		}
	}
}

func (ll *LokiLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	elapsedTime := time.Since(begin)

	sql, rows := fc()

	sqlMessage := fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsedTime.Microseconds())/1000.0, rows, sql)

	switch {
	case err != nil && ll.level >= logger.Error:
		select {
		case logsQueue <- logEntry{"Error executing SQL query:" + sqlMessage, ll.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+sqlMessage, ll.labels))
		}
	case ll.level == logger.Info:
		select {
		case logsQueue <- logEntry{"SQL:" + sqlMessage, ll.labels}:
		default:
			log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+sqlMessage, ll.labels))
		}
	}
}
