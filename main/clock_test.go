package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestClock_printProperMessage(t *testing.T) {
	expectedForSeconds := messages.secMsg + "\n"
	expectedForMinutes := messages.minMsg + "\n"
	expectedForHours := messages.hourMsg + "\n"
	tests := []struct {
		name           string
		elapsedSeconds int
		expected       string
	}{
		{"seconds 1 sec", 1, expectedForSeconds},
		{"seconds 2 secs", 2, expectedForSeconds},
		{"seconds 3 secs", 3, expectedForSeconds},
		{"seconds 1 hour and a second", 3601, expectedForSeconds},
		{"minutes 1 minute", 60, expectedForMinutes},
		{"minutes 2 minutes", 120, expectedForMinutes},
		{"minutes 4 minutes", 240, expectedForMinutes},
		{"hours 1 hour", 3600, expectedForHours},
		{"hours 2 hours", 2 * 3600, expectedForHours},
		{"hours 10 hours", 10 * 3600, expectedForHours},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			buf := bytes.Buffer{}
			c := &Clock{messages: &messages, locking: &testLocking, writer: &buf}
			// act
			c.printProperMessages(tt.elapsedSeconds)
			// assert
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestClock_determineProperMessage(t *testing.T) {
	tests := []struct {
		name            string
		elapsedSeconds  int
		expectedMessage string
	}{
		{"seconds 1 sec", 1, messages.secMsg},
		{"seconds 2 secs", 2, messages.secMsg},
		{"seconds 3 secs", 3, messages.secMsg},
		{"seconds 1 hour and a second", 3601, messages.secMsg},
		{"minutes 1 minute", 60, messages.minMsg},
		{"minutes 2 minutes", 120, messages.minMsg},
		{"minutes 4 minutes", 240, messages.minMsg},
		{"hours 1 hour", 3600, messages.hourMsg},
		{"hours 2 hours", 2 * 3600, messages.hourMsg},
		{"hours 10 hours", 10 * 3600, messages.hourMsg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			c := &Clock{
				messages: &messages,
			}
			// act
			actualMessage := c.determineProperMessage(tt.elapsedSeconds)
			// assert
			assert.Equal(t, tt.expectedMessage, actualMessage)
		})
	}
}

func TestClock_tick(t *testing.T) {
	tickIntervalZero := time.Second * 0
	tests := []struct {
		name           string
		duration       time.Duration
		expectedOutput string
	}{
		{"1 hour", time.Hour, getExpectedOutput(3600, &messages)},
		{"2 hours", time.Hour * 2, getExpectedOutput(3600*2, &messages)},
		{"3 hours", time.Hour * 3, getExpectedOutput(3600*3, &messages)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			buf := &bytes.Buffer{}
			testLocking.wg.Add(1)
			c := &Clock{
				messages:     &messages,
				locking:      &testLocking,
				writer:       buf,
				duration:     tt.duration,
				tickInterval: tickIntervalZero}
			// act
			c.tick()
			// assert
			assert.Equal(t, tt.expectedOutput, buf.String())

		})
	}
}

func getExpectedOutput(seconds int, messages *Messages) string {
	hourMsg := (*messages).hourMsg + "\n"
	minMsg := (*messages).minMsg + "\n"
	secMsg := (*messages).secMsg + "\n"
	var sb strings.Builder
	for i := 1; i <= seconds; i++ {
		if i%3600 == 0 {
			sb.WriteString(hourMsg)
		} else if i%60 == 0 {
			sb.WriteString(minMsg)
		} else {
			sb.WriteString(secMsg)
		}
	}
	return sb.String()
}
