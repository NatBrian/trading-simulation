package helper

import (
	"encoding/json"
	"net/http"
)

// ResponseFormatter returning formatted JSON response
func ResponseFormatter(w http.ResponseWriter, statuscode int, body interface{}) {
	response, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(statuscode)
	_, _ = w.Write(response)
}
