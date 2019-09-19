package main

import (
	"github.com/Sereger/timeoutMutex"
	"log"
	"sync"
	"time"
)

func main() {
	l := timeoutMutex.NewLock(timeoutMutex.WithMaxOverlapReaders(3))
	wg := new(sync.WaitGroup)
	wg.Add(5)
	for i := 0; i < 5; i++ {
		log.Printf("index %d try to take the lock", i)
		l.RLock()
		log.Printf("index %d take the lock", i)
		go func(i int, wg *sync.WaitGroup) {
			time.Sleep(3 * time.Second)
			log.Printf("index %d release the lock", i)
			l.RUnlock()
			wg.Done()
		}(i, wg)
	}

	wg.Wait()
}
