package main

import (
	"github.com/stretchr/testify/assert"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

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
