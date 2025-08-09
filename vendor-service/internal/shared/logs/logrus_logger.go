package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
)

type LokiLogrusHook struct {
	labels map[string]string
}

func NewLokiLogrusLogger(newLabels map[string]string) *LokiLogrusHook {
	return &LokiLogrusHook{labels: newLabels}
}

func (llh *LokiLogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (llh *LokiLogrusHook) Fire(entry *logrus.Entry) error {

	message, err := entry.String()
	if err != nil {
		return err
	}

	labels := make(map[string]string)

	for key, value := range llh.labels {
		labels[key] = value
	}

	labels["level"] = entry.Level.String()

	select {
	case logsQueue <- logEntry{message, labels}:
	default:
		log.Printf("Logs queue is full, dropping log: %s", fmt.Sprintf("ERROR: "+message, llh.labels))

	}

	return nil

}
