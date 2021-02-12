package main

import (
	"os"
	"sync"
	"time"
)

var (
	clockWorkingDuration = time.Hour * 3
	tickInterval         = time.Second
	confFileName         = "file.conf"
	updateConfInterval   = time.Second * 2
	infinity             = time.Hour * 24 * 365 * 290
	writer               = os.Stdout
	messages             = Messages{
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
	clock := &Clock{
		messages:     &messages,
		locking:      &locking,
		writer:       writer,
		duration:     clockWorkingDuration,
		tickInterval: tickInterval,
	}
	confUpdater := &ConfUpdater{
		messages:    &messages,
		locking:     &locking,
		filename:    confFileName,
		interval:    updateConfInterval,
		checkPeriod: infinity,
	}
	locking.wg.Add(2)
	go clock.tick()
	go confUpdater.updateConf()
	locking.wg.Wait()
}
