package main

import (
	"os"
	"sync"
	"time"
)

var (
	clockWorkingDuration       = time.Hour * 3
	tickInterval               = time.Second
	confFileName               = "file.conf"
	updateConfIntervalDuration = time.Second * 2
	infinity                   = time.Hour * 24 * 365 * 290
	messages                   = Messages{
		secMsg:  "tick",
		minMsg:  "tock",
		hourMsg: "bong",
	}
	locking = Locking{
		wg:       sync.WaitGroup{},
		mutex:    sync.Mutex{},
		finished: false,
	}
)

func main() {
	writer := os.Stdout
	locking.wg.Add(1)
	c := &Clock{
		messages:     &messages,
		locking:      &locking,
		writer:       writer,
		duration:     clockWorkingDuration,
		tickInterval: tickInterval,
	}
	go c.tick()
	locking.wg.Add(1)
	go updateConf(&messages, &locking, confFileName, updateConfIntervalDuration, infinity)
	locking.wg.Wait()
}
