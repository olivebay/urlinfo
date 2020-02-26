package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	mgo "gopkg.in/mgo.v2"
)

// FakeDatabase satisfies DataLayer and act as a mock.
type FakeDatabase struct{}

// FakeCollection satisfies Collection and act as a mock.
type FakeCollection struct{}


func (fc FakeCollection) Find(query interface{}) *mgo.Query {
	return nil
}

// C mocks mgo.Database(name).Collection(name).
func (db FakeDatabase) C(name string) Collection {
	return FakeCollection{}
}

// FakeSession satisfies Session and act as a mock of *mgo.session.
type FakeSession struct{}

// NewFakeSession mosck NewSession.
func NewFakeSession() Session {
	return FakeSession{}
}

// Close mocks mgo.Session.Close().
func (fs FakeSession) Close() {}

// Copy mocks mgo.Session.Copy().
// Regarding the context of use, no need to actually Copy the mock.
func (fs FakeSession) Copy() Session {
	return fs
}

// DB mocks mgo.Session.DB().
func (fs FakeSession) DB(name string) DataLayer {
	fakeDataset := FakeDatabase{}
	return fakeDataset
}

// GetBlacklists mocks models.GetBlacklists()
func (db FakeDatabase) GetBlacklists(d string) (URL, error) {
	var u string
	if !strings.HasPrefix(d, "http://") && !strings.HasPrefix(d, "https://") {
		u = "https://" + d
	}
	// parse the url to extract the hostname
	parsedURL, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	h := parsedURL.Hostname()

	var res []URL
	dbContent, err := ioutil.ReadFile("../models/testdata/default_urls.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(dbContent, &res)

	for _, u := range res {
		if u.Domain == h {
			u.Positives = true
			u.Total = len(u.Blacklists)
			u.Url = d
			return u, nil
		}
	}

	return URL{}, ErrURLNotFound
}
