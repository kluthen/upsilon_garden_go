package tools

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// IsAPI Tell whether request requires API reply or not.
func IsAPI(req *http.Request) bool {
	return strings.Contains(req.URL.String(), "/api/")
}

// GetInt parse request to get int value.
func GetInt(req *http.Request, key string) (int, error) {
	vars := mux.Vars(req)
	value, err := strconv.Atoi(vars[key])
	if err != nil {
		log.Printf("Web: requested key: %s , not found in: %s", key, req.URL)
		return 0, errors.New("Invalid key requested")
	}
	return value, nil
}

// GenerateAPIError generate a simple JSON reply with error message provided.
func GenerateAPIError(w http.ResponseWriter, message string) {
	var repm = make(map[string]string)
	repm["error"] = message
	json.NewEncoder(w).Encode(repm)
	w.WriteHeader(400)
}
