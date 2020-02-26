package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"

	"context"
	"github.com/olivebay/urlinfo/api/handlers"
	"github.com/olivebay/urlinfo/api/models"
)

func TestGetURLhandler(t *testing.T) {
	tcs := []struct {
		request string
		url     string
		want    string
		status  int
	}{
		{"/urlinfo/1/nsgs-gov.com:80/test", "nsgs-gov.com:80/test", `{"url":"nsgs-gov.com:80/test","domain":"nsgs-gov.com","positives":true,"total":2,"blacklists":["vault","spamHaus"]}`, 200},
		{"/urlinfo/1/nsgs-gov.com:80", "nsgs-gov.com:80", `{"url":"nsgs-gov.com:80","domain":"nsgs-gov.com","positives":true,"total":2,"blacklists":["vault","spamHaus"]}`, 200},
		{"/urlinfo/1/nsgs-gov.com", "nsgs-gov.com", `{"url":"nsgs-gov.com","domain":"nsgs-gov.com","positives":true,"total":2,"blacklists":["vault","spamHaus"]}`, 200},
		{"/urlinfo/1/bbc.co.uk:443/search?q=test", "bbc.co.uk/search?q=test", `{"message":"url not found"}`, 404},
		{"/urlinfo/1/bbc.co.uk:443", "bbc.co.uk", `{"message":"url not found"}`, 404},
		{"/urlinfo/1/", "", "", 400},
	}

	for _, tc := range tcs {
		req, err := http.NewRequest("GET", tc.request, nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{"url": tc.url})
		req = req.WithContext(context.WithValue(req.Context(), "db", models.FakeSession{}))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetURL)
		handler.ServeHTTP(rr, req)

		//Confirm the response has the right status code
		if status := rr.Code; status != tc.status {
			t.Errorf("GetURL handler returned wrong status code: want %d got %d", tc.status, status)
		}

		got := rr.Body.String()
		got = strings.TrimSpace(got)
		if !cmp.Equal(tc.want, got) {
			t.Errorf(tc.want, got)
		}
	}

}
