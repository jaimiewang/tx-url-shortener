package util

import (
	"encoding/json"
	"net/http"
)

func ParseJsonForm(r *http.Request, i interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		return err
	}

	return nil
}

func WriteJsonResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(i)
	if err != nil {
		return
	}

	_, _ = w.Write(b)
	return
}
