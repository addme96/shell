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
	testMessages := Messages{
		secMsg:  "tick",
		minMsg:  "tock",
		hourMsg: "bong",
	}
	filename := path.Join("test_data", "file.conf")
	interval := time.Second * 0
	checkPeriod := time.Millisecond * 500

	t.Run("clock is working", func(t *testing.T) {
		// arrange
		testLocking = Locking{
			wg:       sync.WaitGroup{},
			mutex:    sync.Mutex{},
			finished: false,
		}
		testLocking.wg.Add(1)

		cu := &ConfUpdater{
			messages:    &testMessages,
			locking:     &testLocking,
			filename:    filename,
			interval:    interval,
			checkPeriod: checkPeriod}
		// act
		cu.updateConf()
		// assert
		assert.Equal(t, "tic", testMessages.secMsg)
		assert.Equal(t, "tac", testMessages.minMsg)
		assert.Equal(t, "toe", testMessages.hourMsg)
		// cleanup
	})
	t.Run("clock has finished", func(t *testing.T) {
		// arrange
		testLocking = Locking{
			wg:       sync.WaitGroup{},
			mutex:    sync.Mutex{},
			finished: true,
		}
		testLocking.wg.Add(1)
		cu := &ConfUpdater{
			messages:    &testMessages,
			locking:     &testLocking,
			filename:    filename,
			interval:    interval,
			checkPeriod: checkPeriod,
		}
		// act
		cu.updateConf()
		// assert
		assert.Equal(t, "tick", testMessages.secMsg)
		assert.Equal(t, "tock", testMessages.minMsg)
		assert.Equal(t, "bong", testMessages.hourMsg)
	})
}

func TestConfUpdater_hasClockFinished(t *testing.T) {
	t.Run("clock is working", func(t *testing.T) {
		// arrange
		testLocking = Locking{
			wg:       sync.WaitGroup{},
			mutex:    sync.Mutex{},
			finished: true,
		}
		cu := &ConfUpdater{
			locking: &testLocking,
		}
		// act
		finished := cu.hasClockFinished()
		// assert
		assert.Equal(t, true, finished)
	})
	t.Run("clock has finished", func(t *testing.T) {
		// arrange
		testLocking = Locking{
			wg:       sync.WaitGroup{},
			mutex:    sync.Mutex{},
			finished: false,
		}
		cu := &ConfUpdater{
			locking: &testLocking,
		}
		// act
		finished := cu.hasClockFinished()
		// assert
		assert.Equal(t, false, finished)
	})
}
