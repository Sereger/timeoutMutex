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
