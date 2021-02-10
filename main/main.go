package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Messages struct {
	secMsg  string
	minMsg  string
	hourMsg string
}

type Locking struct {
	wg       sync.WaitGroup
	mutex    sync.Mutex
	finished bool
}

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
	go tickClock(&messages, &locking, writer, clockWorkingDuration, tickInterval)
	locking.wg.Add(1)
	go updateConf(&messages, &locking, confFileName, updateConfIntervalDuration, infinity)
	locking.wg.Wait()
}

func tickClock(
	messages *Messages,
	locking *Locking,
	writer io.Writer,
	duration time.Duration,
	tickInterval time.Duration) {
	durationsSeconds := int(duration / time.Second)
	for i := 1; i <= durationsSeconds; i++ {

		time.Sleep(tickInterval)
		locking.mutex.Lock()
		printProperMessage(messages, i, writer)
		locking.mutex.Unlock()
	}
	(*locking).mutex.Lock()
	(*locking).finished = true
	(*locking).mutex.Unlock()
	(*locking).wg.Done()

}

func printProperMessage(messages *Messages, i int, writer io.Writer) {
	if i%3600 == 0 {
		_, _ = fmt.Fprintln(writer, (*messages).hourMsg)
	} else if i%60 == 0 {
		_, _ = fmt.Fprintln(writer, (*messages).minMsg)
	} else {
		_, _ = fmt.Fprintln(writer, (*messages).secMsg)
	}
}

func updateConf(
	messages *Messages,
	locking *Locking,
	filename string,
	interval time.Duration,
	checkPeriod time.Duration) {
	stopTime := time.Now().Local().Add(checkPeriod)

	for {
		now := time.Now()
		if now.After(stopTime) {
			break
		}
		(*locking).mutex.Lock()
		clockFinished := (*locking).finished
		(*locking).mutex.Unlock()
		if clockFinished {
			break
		}
		time.Sleep(interval)

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		newMessages := readConfig(file)
		_ = file.Close()
		(*locking).mutex.Lock()
		*messages = newMessages
		(*locking).mutex.Unlock()
	}
	(*locking).wg.Done()
}

func readConfig(reader io.Reader) Messages {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	secM := scanner.Text()
	scanner.Scan()
	minM := scanner.Text()
	scanner.Scan()
	hourM := scanner.Text()
	return Messages{
		secMsg:  secM,
		minMsg:  minM,
		hourMsg: hourM,
	}
}
