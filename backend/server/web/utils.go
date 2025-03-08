package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func logAndWriteError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func responseWithJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logAndWriteError(w, err)
	}
}

func readFromJSON(r *http.Request, v any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
