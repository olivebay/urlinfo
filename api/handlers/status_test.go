package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olivebay/urlinfo/api/handlers"
)

func TestStatusHanlder(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.StatusHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
