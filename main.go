package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SessionFilter struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type RequestPayload struct {
	Filters []SessionFilter `json:"filters"`
}

func GetBotSessions(apiToken, startTime, endTime string) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://app.jaicp.com/api/reporter/p/%s/sessions/filter", apiToken)

	// создание request
	// тут type должен быть 'DATE_TIME_RANGE' по условию. но не понятно было что насчет key
	filters := []SessionFilter{
		{
			Key:  "SESSION_START_TIME",
			Type: "DATE_TIME_RANGE",
			From: startTime,
			To:   endTime,
		},
	}
	payload := RequestPayload{Filters: filters}

	// перевод в bytes
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, fmt.Errorf("Error: %d - %s", resp.StatusCode, resp.Status)
}

func main() {
	apiToken := "wKYIZxYP:75acc1d688407813562794aa57bed1f723064c64"
	startTime := "2023-03-15T10:00:00Z"
	endTime := "2023-03-15T12:59:59Z"

	result, err := GetBotSessions(apiToken, startTime, endTime)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// добавление отступов
	resultJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Printf("marshaling err: %v\n", err)
		return
	}

	// запись в файл
	err = os.WriteFile("out.json", resultJSON, 0644)
	if err != nil {
		fmt.Printf("writeFile err: %v\n", err)
		return
	}

	fmt.Println("The result was successfully written to out.json")

}
