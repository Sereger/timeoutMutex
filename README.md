# timeoutMutex

A mutex is a reader/writer mutual exclusion lock.
The lock can be held by an some number of readers or a single writer.
The lock returned an error when the waiting time will exceed specified limit.

Examlpe:
```go
package main

import (
	"github.com/Sereger/timeoutMutex"
	"log"
	"time"
)

func main() {
	l := timeoutMutex.NewLock()
	l.Lock()
	log.Println("Lock")

	time.AfterFunc(5*time.Second, func() {
		log.Println("Unlock after 5 sec")
		l.Unlock()
	})

	log.Println("try to take lock with timeout 1 sec")
	err := l.TimeoutLock(time.Second)
	if err != nil {
		log.Println("receive err:", err)
	} else {
		// we don't should be here!!!
		log.Fatal("took the lock!")
	}

	log.Println("try to take lock with timeout 5 sec")
	err = l.TimeoutLock(5 * time.Second)
	if err != nil {
		log.Fatal("receive err:", err)
	} else {
		log.Println("took the lock!")
	}
}
```

Output:
```bash
00:20:56 Lock
00:20:56 try to take lock with timeout 1 sec
00:20:57 receive err: timeout
00:20:57 try to take lock with timeout 5 sec
00:21:01 Unlock after 5 sec
00:21:01 took lock!
```


#### Reades limit
Sometimes, you need to limit the number of readers, so you can use `WithMaxOverlapReaders` option for this. 
Example:
```go
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

```

Output
```bash
00:00:00 index 0 try to take the lock
00:00:00 index 0 take the lock
00:00:00 index 1 try to take the lock
00:00:00 index 1 take the lock
00:00:00 index 2 try to take the lock
00:00:00 index 2 take the lock
00:00:00 index 3 try to take the lock
00:00:03 index 0 release the lock
00:00:03 index 2 release the lock
00:00:03 index 3 take the lock
00:00:03 index 4 try to take the lock
00:00:03 index 4 take the lock
00:00:03 index 1 release the lock
00:00:06 index 4 release the lock
00:00:06 index 3 release the lock
```