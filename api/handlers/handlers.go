package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olivebay/urlinfo/api/models"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// GetURL handler checks if the url is in the blacklist datastore and
// return JSON containing the blacklisted domains
func GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := vars["url"]

	w.Header().Set("Content-Type", "application/json")

	// return 404 if a url is not provided
	if u == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := MgoDBFromR(r)

	//fetch the url from the datastore
	urlinfo, err := db.GetBlacklists(u)
	if err == models.ErrURLNotFound {
		w.WriteHeader(http.StatusNotFound)
		models.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	// serialize urlinfo struct into JSON
	err = models.ToJSON(urlinfo, w)
	if err != nil {
		log.Println("[ERROR] serializing Url", err)
	}
	return
}

// MgoDBFromR takes a request argument and return the extracted *mgo.session.
func MgoDBFromR(r *http.Request) models.DataLayer {
	return MgoSessionFromCtx(r.Context()).DB("urls")
}

// MgoSessionFromCtx takes a context argument and return the related *mgo.session.
func MgoSessionFromCtx(ctx context.Context) models.Session {
	mgoSession, _ := ctx.Value("db").(models.Session)
	return mgoSession
}
