package logs

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

type Payload struct {
	Streams []LokiStream `json:"streams"`
}

func SendLogsToLoki(message string, labels map[string]string) error {

	stackId := os.Getenv("GRAFANA_LOKI_STACK_ID")
	apiKey := os.Getenv("GRAFANA_LOKI_API_KEY")
	url := os.Getenv("GRAFANA_LOKI_URL")

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	stream := LokiStream{
		Stream: labels,
		Values: [][2]string{{timestamp, message}},
	}

	payload := Payload{
		Streams: []LokiStream{stream},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	auth := stackId + ":" + apiKey
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	request.Header.Set("Authorization", "Basic "+encodedAuth)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			return
		}
	}()

	if response.StatusCode >= 300 {
		return fmt.Errorf("error sending logs to Loki, status: %s", response.Status)
	}

	return nil
}
