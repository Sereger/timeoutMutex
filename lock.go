package timeoutMutex

import (
	"errors"
	"math"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	unlocked = uint32(0)
	locked   = uint32(1)
)

var Timeout = errors.New("timeout")

type (
	RWLock struct {
		buff             [8]byte
		maxR             uint32
		addrLock         *uint32
		addrCnr          *uint32
		sleepingDuration time.Duration
	}

	Option func(l *RWLock)
)

func WithMaxOverlapReaders(m uint32) Option {
	return func(l *RWLock) {
		l.maxR = m
	}
}

func WithSleepingTimeout(d time.Duration) Option {
	return func(l *RWLock) {
		l.sleepingDuration = d
	}
}

func NewLock(opts ...Option) *RWLock {
	b := [8]byte{}

	l := &RWLock{
		buff:             b,
		maxR:             math.MaxUint32,
		addrLock:         (*uint32)(unsafe.Pointer(&b)),
		addrCnr:          (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 4)),
		sleepingDuration: 500 * time.Nanosecond,
	}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *RWLock) TimeoutLock(d time.Duration) error {
	s := time.Now()
	for !atomic.CompareAndSwapUint32(l.addrLock, unlocked, locked) {
		if time.Since(s) > d {
			return Timeout
		}
		l.sleep()
	}

	for atomic.LoadUint32(l.addrCnr) != 0 {
		if time.Since(s) > d {
			l.Unlock()
			return Timeout
		}
		l.sleep()
	}

	return nil
}
func (l *RWLock) TimeoutRLock(d time.Duration) error {
	s := time.Now()
	for {
		if time.Since(s) > d {
			return Timeout
		}
		if l.locked() {
			l.sleep()
			continue
		}
		nv := l.rIncr()
		if nv > l.maxR {
			l.rDecr()
			l.sleep()
			continue
		}

		break
	}

	return nil
}
func (l *RWLock) Lock() {
	for !atomic.CompareAndSwapUint32(l.addrLock, unlocked, locked) {
		l.sleep()
	}
	for atomic.LoadUint32(l.addrCnr) != 0 {
		l.sleep()
	}
}

func (l *RWLock) RLock() {
	for {
		if l.locked() {
			l.sleep()
			continue
		}

		nv := l.rIncr()
		if nv > l.maxR {
			l.rDecr()
			l.sleep()
			continue
		}

		break
	}
}
func (l *RWLock) locked() bool {
	return atomic.LoadUint32(l.addrLock) == locked
}
func (l *RWLock) RUnlock() {
	l.rDecr()
}

func (l *RWLock) rIncr() uint32 {
	return atomic.AddUint32(l.addrCnr, 1)
}

func (l *RWLock) rDecr() uint32 {
	return atomic.AddUint32(l.addrCnr, ^uint32(0))
}

func (l *RWLock) Unlock() {
	atomic.StoreUint32(l.addrLock, 0)
}

func (l *RWLock) sleep() {
	time.Sleep(l.sleepingDuration)
}
