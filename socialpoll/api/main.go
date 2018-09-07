package main

import (
	"flag"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "Address for Endpoint")
		mongo = flag.String("mongo", "localhost", "Address for MongoDB")
	)
	flag.Parse()

	log.Println("Dialing MongoDB... ", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("Failed to Connect to MongoDB")
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))

	log.Println("Starting web server on ", *addr)
	http.ListenAndServe(":8080", mux)
	log.Println("Stopping...")
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

func withData(d *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()

		SetVar(r, "db", thisDb.DB("ballots"))
		f(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)

		fn(w, r)
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")

		fn(w, r)
	}
}
