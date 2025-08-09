package logs

import "log"

type logEntry struct {
	message string
	labels  map[string]string
}

var logsQueue = make(chan logEntry, 4000)

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
