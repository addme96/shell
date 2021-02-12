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
	testLocking = Locking{
		wg:       sync.WaitGroup{},
		mutex:    sync.Mutex{},
		finished: false,
	}
)

func main() {
	clock := &Clock{
		messages:     &messages,
		locking:      &testLocking,
		writer:       writer,
		duration:     clockWorkingDuration,
		tickInterval: tickInterval,
	}
	confUpdater := &ConfUpdater{
		messages:    &messages,
		locking:     &testLocking,
		filename:    confFileName,
		interval:    updateConfInterval,
		checkPeriod: infinity,
	}
	testLocking.wg.Add(2)
	go clock.tick()
	go confUpdater.updateConf()
	testLocking.wg.Wait()
}
