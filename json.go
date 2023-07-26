package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func writeJson(w http.ResponseWriter, status int, payload interface{}, headers ...http.Header) error {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func errorJSON(w http.ResponseWriter, msg string, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		log.Printf("Responding with %v: %v", status[0], msg)
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = msg

	return writeJson(w, statusCode, payload)

}
