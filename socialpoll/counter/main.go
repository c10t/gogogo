package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("Connect to MongoDB...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}

	defer func() {
		log.Println("Close connection for MongoDB...")
		db.Close()
	}()

	// pollData := db.DB("ballots").C("polls")
}