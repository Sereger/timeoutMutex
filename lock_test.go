package timeoutMutex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLock1(t *testing.T) {
	lockLimit := 10
	l := NewLock(WithMaxOverlapReaders(uint32(lockLimit)))
	wg := new(sync.WaitGroup)
	moment := time.Now()
	sleepRLock := time.Second * 2
	sleepLock := time.Second * 3
	var ok, fail uint32

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int, ok *uint32, fail *uint32, lk *RWLock, wg *sync.WaitGroup) {
			defer wg.Done()
			err := lk.TimeoutRLock(time.Second * 10)
			if err != nil {
				atomic.AddUint32(fail, 1)
				t.Logf("iter %d: %s", i, err)
				return
			}
			defer lk.RUnlock()
			s := time.Now()
			t.Logf("iter %d started", i)
			time.Sleep(sleepRLock)
			t.Logf("iter %d ended (%s)", i, time.Since(s))
			atomic.AddUint32(ok, 1)
		}(i, &ok, &fail, l, wg)
	}

	time.AfterFunc(time.Second*3, func() {
		t.Logf("Lock enter")
		l.Lock()
		defer l.Unlock()
		t.Logf("Lock start")
		time.Sleep(sleepLock)
		t.Logf("Lock end")
	})

	wg.Wait()
	tDiff := time.Since(moment).Round(20 * time.Millisecond)
	expDuration := 2*sleepRLock + sleepLock + 2*sleepRLock
	if tDiff != expDuration {
		t.Fatalf("incorrect execution time: [%s], expect: [%s]", tDiff, expDuration)
	}

	t.Logf("%s, ok: %d, f: %d", time.Since(moment), ok, fail)
	if ok != 40 {
		t.Fail()
	}
}
func TestLock2(t *testing.T) {
	l := NewLock()
	wg := new(sync.WaitGroup)
	s := time.Now()
	var ok, fail uint32
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int, ok *uint32, fail *uint32, lk *RWLock, wg *sync.WaitGroup) {
			defer wg.Done()
			err := lk.TimeoutRLock(time.Second * 10)
			if err != nil {
				atomic.AddUint32(fail, 1)
				t.Logf("iter %d: %s", i, err)
				return
			}
			defer lk.RUnlock()
			s := time.Now()
			t.Logf("iter %d started", i)
			time.Sleep(time.Second * 2)
			t.Logf("iter %d ended (%s)", i, time.Since(s))
			atomic.AddUint32(ok, 1)
		}(i, &ok, &fail, l, wg)
	}

	wg.Wait()
	tDiff := time.Since(s)

	t.Logf("%s, ok: %d, f: %d", tDiff.Round(time.Millisecond), ok, fail)
	if ok != 200 {
		t.Fail()
	}
}

func TestLock3(t *testing.T) {
	l := NewLock()
	wg := new(sync.WaitGroup)
	moment := time.Now()
	sleepDuration := time.Second * 2
	var ok, fail uint32
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int, ok *uint32, fail *uint32, lk *RWLock, wg *sync.WaitGroup) {
			defer wg.Done()
			err := lk.TimeoutRLock(time.Second * 10)
			if err != nil {
				atomic.AddUint32(fail, 1)
				t.Logf("iter %d: %s", i, err)
				return
			}
			defer lk.RUnlock()
			time.Sleep(sleepDuration)
			atomic.AddUint32(ok, 1)
		}(i, &ok, &fail, l, wg)
	}

	wg.Wait()
	tDiff := time.Since(moment).Round(50 * time.Millisecond)
	if tDiff != sleepDuration {
		t.Fatalf("incorrect execution time: [%s], expect: [%s]", tDiff, sleepDuration)
	}

	t.Logf("%s, ok: %d, f: %d", tDiff.Round(time.Millisecond), ok, fail)
	if ok != 200 {
		t.Fatalf("ok is: %d", ok)
	}
}
