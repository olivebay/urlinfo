package models

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// ErrURLNotFound is returned when an url is not found
var ErrURLNotFound = fmt.Errorf("url not found")

// URL represents information about a URL
type URL struct {
	Url   string   `json:"url"`  
	Domain     string   `json:"domain"`   
	Positives  bool     `json:"positives"` 
	Total      int      `json:"total"`
	Blacklists []string `json:"blacklists"`
}

// GetBlacklists returns a list of blacklists that a domain was found
func (db *MongoDatabase) GetBlacklists(d string) (URL, error) {
	// parse the url to extract the hostname
	var u string
	if !strings.HasPrefix(d, "http://") && !strings.HasPrefix(d, "https://") {
		u = "https://" + d
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	h := parsedURL.Hostname()

	// fetch the domain from the datastore
	var result URL
	err = db.C(collection).Find(bson.M{"domain": h}).One(&result)
	if err != nil {
		return URL{}, ErrURLNotFound
	}
	result.Positives = true
	result.Url = d
	result.Total = len(result.Blacklists)

	return result, nil
}
