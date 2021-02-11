package main

import (
	"fmt"
	"io"
	"time"
)

type Clock struct {
	messages     *Messages
	locking      *Locking
	writer       io.Writer
	duration     time.Duration
	tickInterval time.Duration
}

func (c *Clock) tickClock() {
	durationsSeconds := int(c.duration / time.Second)
	for i := 1; i <= durationsSeconds; i++ {
		time.Sleep(c.tickInterval)
		c.printProperMessages(i)
	}
	(*c.locking).mutex.Lock()
	(*c.locking).finished = true
	(*c.locking).mutex.Unlock()
	(*c.locking).wg.Done()
}

func (c *Clock) printProperMessages(secondsElapsed int) {
	c.locking.mutex.Lock()
	msg := c.determineProperMessage(secondsElapsed)
	c.locking.mutex.Unlock()
	c.writeString(msg)
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

func (c *Clock) writeString(s string) {
	_, _ = fmt.Fprintln(c.writer, s)
}
