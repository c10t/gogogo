package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mgo "gopkg.in/mgo.v2"
)

// Session is an interface to access to the Session struct.
type Session interface {
	DB(name string) DataLayer
	Close()
}

// MongoSession is currently a Mongo session.
type MongoSession struct {
	*mgo.Session
}

// DB shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (s MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

// DataLayer is an interface to access to the database struct.
type DataLayer interface {
	C(name string) Collection
}

// MockSession satisfies Session and act as a mock of *mgo.session.
type MockSession struct{}

// NewMockSession mock NewSession.
func NewMockSession() Session {
	return MockSession{}
}

// Close mocks mgo.Session.Close().
func (fs MockSession) Close() {}

// DB mocks mgo.Session.DB().
func (fs MockSession) DB(name string) DataLayer {
	mockDatabase := MockDatabase{}
	return mockDatabase
}

// NewSession returns a new Mongo Session.
func NewSession() Session {
	mgoSession, err := mgo.Dial("<MONGO_URI>")
	if err != nil {
		panic(err)
	}
	return MongoSession{mgoSession}
}

// MongoCollection wraps a mgo.Collection to embed methods in models.
type MongoCollection struct {
	*mgo.Collection
}

// Collection is an interface to access to the collection struct.
type Collection interface {
	Find(query interface{}) *mgo.Query
	Count() (n int, err error)
	Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	Update(selector interface{}, update interface{}) error
}

// MongoDatabase wraps a mgo.Database to embed methods in models.
type MongoDatabase struct {
	*mgo.Database
}

// C shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// MockDatabase satisfies DataLayer and act as a mock.
type MockDatabase struct{}

// MockCollection satisfies Collection and act as a mock.
type MockCollection struct{}

// Find mock.
func (fc MockCollection) Find(query interface{}) *mgo.Query {
	return nil
}

// Count mock.
func (fc MockCollection) Count() (n int, err error) {
	return 10, nil
}

// Insert mock.
func (fc MockCollection) Insert(docs ...interface{}) error {
	return nil
}

// Remove mock.
func (fc MockCollection) Remove(selector interface{}) error {
	return nil
}

// Update mock.
func (fc MockCollection) Update(selector interface{}, update interface{}) error {
	return nil
}

// C mocks mgo.Database(name).Collection(name).
func (db MockDatabase) C(name string) Collection {
	return MockCollection{}
}

func TestGet(t *testing.T) {
	t.Logf("vars state is: %v", vars)

	mockServer := httptest.NewServer(http.HandlerFunc(handlePolls))
	defer mockServer.Close()

	r, err := http.Get(mockServer.URL + "/people/")

	if err != nil {
		t.Logf("mocking now: %v, %v", r, err)
	}
}
