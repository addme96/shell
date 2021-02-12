package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

type ConfUpdater struct {
	messages    *Messages
	locking     *Locking
	filename    string
	interval    time.Duration
	checkPeriod time.Duration
}

func (c *ConfUpdater) updateConf() {
	stopTime := time.Now().Local().Add(c.checkPeriod)
	for !time.Now().After(stopTime) && !c.hasClockFinished() {
		time.Sleep(c.interval)
		file, err := os.Open(c.filename)
		if err != nil {
			log.Fatal(err)
		}
		newMessages := c.readConfig(file)
		_ = file.Close()
		(*c.locking).mutex.Lock()
		*c.messages = newMessages
		(*c.locking).mutex.Unlock()
	}
	(*c.locking).wg.Done()
}

func (c *ConfUpdater) hasClockFinished() bool {
	(*c.locking).mutex.Lock()
	defer (*c.locking).mutex.Unlock()
	clockFinished := (*c.locking).finished
	return clockFinished
}

func (c *ConfUpdater) readConfig(reader io.Reader) Messages {
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
