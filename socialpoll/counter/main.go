package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const updateDuration = 1 * time.Second

func main() {
	err := counterMain()
	if err != nil {
		log.Fatal(err)
	}
}

func counterMain() error {
	log.Println("Connect to MongoDB...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		return err
	}

	defer func() {
		log.Println("Close connection for MongoDB...")
		db.Close()
	}()

	pollData := db.DB("ballots").C("polls")

	var countsLock sync.Mutex
	var counts map[string]int

	log.Println("Connect to NSQ...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		return err
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()

		if counts == nil {
			counts = make(map[string]int)
		}

		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		return err
	}

	log.Println("waiting for vote on NSQ...")

	ticker := time.NewTicker(updateDuration)
	defer ticker.Stop()

	update := func() {
		countsLock.Lock()
		defer countsLock.Unlock()

		if len(counts) == 0 {
			log.Println("Skip update DB because there is no vote")
			return
		}

		log.Println("Update DB...")
		log.Println(counts)
		ok := true

		for option, count := range counts {
			sel := bson.M{"options": bson.M{"$in": []string{option}}}
			up := bson.M{"$inc": bson.M{"results." + option: count}}

			if _, err := pollData.UpdateAll(sel, up); err != nil {
				log.Println("failed to update:", err)
				ok = false
			} else {
				counts[option] = 0
			}
		}

		if ok {
			log.Println("Finish to update DB")
			counts = nil
		}
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-ticker.C:
			update()
		case <-termChan:
			q.Stop()
		case <-q.StopChan:
			return nil // finished
		}
	}
}
