package database

import "sync"

var (
	currentSequenceValue = 0
	currentSequenceLock  sync.Mutex
)

func NextSequenceValue() int {
	currentSequenceLock.Lock()
	defer currentSequenceLock.Unlock()

	v := currentSequenceValue
	currentSequenceValue++
	return v
}
