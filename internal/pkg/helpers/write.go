package helpers

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func WriteXML(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	err := xml.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
