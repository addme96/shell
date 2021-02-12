package main

import (
	"github.com/stretchr/testify/assert"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConfUpdater_readConfig(t *testing.T) {
	cu := &ConfUpdater{}
	messages2 := Messages{
		secMsg:  "tic",
		minMsg:  "tac",
		hourMsg: "toe",
	}
	tests := []struct {
		name             string
		configToRead     string
		expectedMessages Messages
	}{
		{"tic tac toe", "tic\ntac\ntoe\n", messages2},
		{"tick tock bong", "tick\ntock\nbong\n", messages},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			reader := strings.NewReader(tt.configToRead)
			// act
			actualMessages := cu.readConfig(reader)
			// assert
			assert.Equal(t, tt.expectedMessages.secMsg, actualMessages.secMsg)
			assert.Equal(t, tt.expectedMessages.minMsg, actualMessages.minMsg)
			assert.Equal(t, tt.expectedMessages.hourMsg, actualMessages.hourMsg)
		})
	}
}

func TestConfUpdater_updateConf(t *testing.T) {
	t.Run("clock is working", func(t *testing.T) {
		// arrange
		oldMessages := messages
		locking = Locking{
			wg:       sync.WaitGroup{},
			mutex:    sync.Mutex{},
			finished: false,
		}
		locking.wg.Add(1)
		filename := path.Join("test_data", "file.conf")
		interval := time.Second * 0
		checkPeriod := time.Millisecond * 500
		cu := &ConfUpdater{
			messages:    &messages,
			locking:     &locking,
			filename:    filename,
			interval:    interval,
			checkPeriod: checkPeriod}
		// act
		cu.updateConf()
		// assert
		assert.Equal(t, "tic", messages.secMsg)
		assert.Equal(t, "tac", messages.minMsg)
		assert.Equal(t, "toe", messages.hourMsg)

		// cleanup
		messages = oldMessages
	})
	t.Run("clock has finished", func(t *testing.T) {
		// arrange
		locking = Locking{
			wg:       sync.WaitGroup{},
			mutex:    sync.Mutex{},
			finished: true,
		}
		locking.wg.Add(1)
		filename := path.Join("test_data", "file.conf")
		interval := time.Second * 0
		checkPeriod := time.Millisecond * 500
		cu := &ConfUpdater{
			messages:    &messages,
			locking:     &locking,
			filename:    filename,
			interval:    interval,
			checkPeriod: checkPeriod,
		}
		// act
		cu.updateConf()
		// assert
		assert.Equal(t, "tick", messages.secMsg)
		assert.Equal(t, "tock", messages.minMsg)
		assert.Equal(t, "bong", messages.hourMsg)
	})
}
