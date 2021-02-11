package main

import "sync"

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
