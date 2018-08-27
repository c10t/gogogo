package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

func main() {
	if err := dialdb(); err != nil {
		log.Fatalln("Failed to dial MongoDB:", err)
	}
	defer closedb()

	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("main: stopping...")
		stopChan <- struct{}{}
		closeConn()
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
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

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())

	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}

		log.Println("Publisher: Stopping...")
		pub.Stop()
		log.Println("Publisher: Stopped.")
		stopchan <- struct{}{}
	}()

	return stopchan
}
