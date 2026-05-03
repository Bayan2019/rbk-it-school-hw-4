package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// _ = json.NewEncoder(w).Encode(v)
	dat, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error setting Message is set: %s", err)
	}
}

func WriteError(w http.ResponseWriter, code int, msg string, err error) {
	// if err != nil {
	// 	log.Println(err)
	// }
	// if code > 499 {
	// 	log.Printf("Responding with 5XX error: %s", msg)
	// }

	WriteJSON(w, code, errorResponse{
		Error:   fmt.Sprintf("%s: %v", msg, err),
		Message: msg,
	})
}
