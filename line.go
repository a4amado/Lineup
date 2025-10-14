package line

import (
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	line          []*QueueItem
	lock          sync.Mutex
	maxProcessing int
}

type QueueItem struct {
	Halt   chan int
	purged bool
	mutex  sync.Mutex
}

func (qi *QueueItem) Purge() {
	qi.mutex.Lock()
	defer qi.mutex.Unlock()
	qi.purged = true
}

func (q *Queue) recalibrateIndexes() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC in recalibrateIndexes:", r)
		}
	}()

	for {

		newLine := []*QueueItem{}
		q.lock.Lock()
		fmt.Println(q.line)
		for _, item := range q.line {

			if item.purged {
				continue
			}

			newLine = append(newLine, item)

		}

		for idx, item := range newLine {

			if idx < q.maxProcessing+1 {
				select {
				case item.Halt <- 1:
				default:
				}
			} else {
				select {
				case item.Halt <- idx - q.maxProcessing + 1:
				default:
				}

			}

		}

		q.line = newLine
		q.lock.Unlock()
		time.Sleep(300 * time.Millisecond)
	}
}

type QueueOptions struct {
	MaxProcessing int
}

func New(opts QueueOptions) *Queue {
	q := &Queue{
		line:          []*QueueItem{},
		lock:          sync.Mutex{},
		maxProcessing: opts.MaxProcessing,
	}

	go q.recalibrateIndexes()
	return q
}

func (q *Queue) Place() *QueueItem {

	newChan := make(chan int, 1)
	newItem := &QueueItem{
		Halt:   newChan,
		purged: false,
	}
	q.lock.Lock()
	q.line = append(q.line, newItem)
	q.lock.Unlock()
	return newItem
}
