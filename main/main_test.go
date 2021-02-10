package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_printProperMessage(t *testing.T) {
	messages := &Messages{
		secMsg:  "tick",
		minMsg:  "tock",
		hourMsg: "bong",
	}
	type args struct {
		messages *Messages
		i        int
	}

	expectedForSeconds := messages.secMsg + "\n"
	expectedForMinutes := messages.minMsg + "\n"
	expectedForHours := messages.hourMsg + "\n"
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{"seconds 1 sec", args{messages, 1}, expectedForSeconds},
		{"seconds 2 secs", args{messages, 2}, expectedForSeconds},
		{"seconds 3 secs", args{messages, 3}, expectedForSeconds},
		{"seconds 1 hour and a second", args{messages, 3601}, expectedForSeconds},
		{"minutes 1 minute", args{messages, 60}, expectedForMinutes},
		{"minutes 2 minutes", args{messages, 120}, expectedForMinutes},
		{"minutes 4 minutes", args{messages, 240}, expectedForMinutes},
		{"hours 1 hour", args{messages, 3600}, expectedForHours},
		{"hours 2 hours", args{messages, 2 * 3600}, expectedForHours},
		{"hours 10 hours", args{messages, 10 * 3600}, expectedForHours},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			buf := bytes.Buffer{}
			// act
			printProperMessage(tt.args.messages, tt.args.i, &buf)
			// assert
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func Test_readConfig(t *testing.T) {
	messages1 := Messages{
		secMsg:  "tic",
		minMsg:  "tac",
		hourMsg: "toe",
	}
	messages2 := Messages{
		secMsg:  "tick",
		minMsg:  "tock",
		hourMsg: "bong",
	}

	tests := []struct {
		name             string
		config           string
		expectedMessages Messages
	}{
		{"tic tac toe", "tic\ntac\ntoe\n", messages1},
		{"tick tock bong", "tick\ntock\nbong\n", messages2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			reader := strings.NewReader(tt.config)
			// act
			actualMessages := readConfig(reader)
			// assert
			assert.Equal(t, tt.expectedMessages.secMsg, actualMessages.secMsg)
			assert.Equal(t, tt.expectedMessages.minMsg, actualMessages.minMsg)
			assert.Equal(t, tt.expectedMessages.hourMsg, actualMessages.hourMsg)
		})
	}
}

func Test_tickClock(t *testing.T) {
	tickInterval := time.Second * 0
	messages = Messages{
		secMsg:  "tick",
		minMsg:  "tock",
		hourMsg: "bong",
	}
	locking = Locking{
		wg:       sync.WaitGroup{},
		mutex:    sync.Mutex{},
		finished: false,
	}

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
			locking.wg.Add(1)

			// act
			tickClock(&messages, &locking, buf, tt.duration, tickInterval)

			// assert
			assert.Equal(t, tt.expectedOutput, buf.String())

		})
	}
}

func getExpectedOutput(seconds int, messages *Messages) string {
	var sb strings.Builder
	hourMsg := (*messages).hourMsg + "\n"
	minMsg := (*messages).minMsg + "\n"
	secMsg := (*messages).secMsg + "\n"
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

func Test_updateConf(t *testing.T) {
	// arrange
	messages = Messages{
		secMsg:  "tick",
		minMsg:  "tock",
		hourMsg: "bong",
	}
	locking = Locking{
		wg:       sync.WaitGroup{},
		mutex:    sync.Mutex{},
		finished: false,
	}
	filename := path.Join("test_data", "file.conf")
	interval := time.Second * 0
	checkPeriod := time.Millisecond * 500
	locking.wg.Add(1)

	// act
	updateConf(&messages, &locking, filename, interval, checkPeriod)
	// assert
	assert.Equal(t, "tic", messages.secMsg)
	assert.Equal(t, "tac", messages.minMsg)
	assert.Equal(t, "toe", messages.hourMsg)
}
