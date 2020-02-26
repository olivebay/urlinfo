package handlers

import (
	"net/http"
)

// StatusHandler endpoint to acknowledge application status.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}