package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type LokiHook struct {
	Labels map[string]string
}

func (h *LokiHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LokiHook) Fire(entry *logrus.Entry) error {

	msg, err := entry.String()
	if err != nil {
		return err
	}

	labels := make(map[string]string)
	for key, value := range h.Labels {
		labels[key] = value
	}

	labels["level"] = entry.Level.String()

	go func() {
		if err := SendLogsToLoki(msg, labels); err != nil {
			fmt.Println("Loki logrus error:", err)
		}
	}()

	return nil
}
