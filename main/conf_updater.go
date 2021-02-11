package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

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
