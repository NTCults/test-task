package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ResponseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	ResponseJSON(w, code, map[string]string{"error": message})
}

func Getenv(key string, defaultVal int) int {
	value := os.Getenv(key)

	if len(value) == 0 {
		return defaultVal
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Wrong environment variable value: '%s:%s'", key, value)
	}
	return intValue
}
