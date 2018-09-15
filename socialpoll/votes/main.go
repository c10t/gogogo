package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		cancel()
		log.Println("main: stopping...")
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("Failed to dial MongoDB:", err)
	}
	defer closedb()

	// start the process
	votes := make(chan string)
	go twitterStream(ctx, votes)
	publishVotes(votes)
}

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("Connecting to database... -> localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("Database connection closed.")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()

	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}

	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) {
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	for vote := range votes {
		pub.Publish("votes", []byte(vote))
	}
	log.Println("Publisher: stopping...")
	pub.Stop()
	log.Println("Publisher: stopped")
}
