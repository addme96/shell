package main

import (
	"fmt"
	"io"
	"log"
	"time"
)

type Clock struct {
	messages     *Messages
	locking      *Locking
	writer       io.Writer
	duration     time.Duration
	tickInterval time.Duration
}

func (c *Clock) tick() {
	defer (*c.locking).wg.Done()
	durationsSeconds := int(c.duration / time.Second)
	for i := 1; i <= durationsSeconds; i++ {
		time.Sleep(c.tickInterval)
		if err := c.printProperMessages(i); err != nil {
			log.Fatal("Failed to print ticker message due to the following error: " + err.Error())
		}
	}
	(*c.locking).mutex.Lock()
	defer (*c.locking).mutex.Unlock()
	(*c.locking).finished = true
}

func (c *Clock) printProperMessages(secondsElapsed int) error {
	c.locking.mutex.Lock()
	defer c.locking.mutex.Unlock()
	msg := c.determineProperMessage(secondsElapsed)
	_, err := fmt.Fprintln(c.writer, msg)
	return err
}

func (c *Clock) determineProperMessage(secondsElapsed int) string {
	if secondsElapsed%3600 == 0 {
		return (*c.messages).hourMsg
	} else if secondsElapsed%60 == 0 {
		return (*c.messages).minMsg
	} else {
		return (*c.messages).secMsg
	}
}
