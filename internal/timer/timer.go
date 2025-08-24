package timer

import (
	"sync"
	"time"
)

type Timer struct {
	val    byte
	action func()
	mu     sync.Mutex
}

func NewTimer() *Timer {
	timer := &Timer{}
	go timer.run()
	return timer
}

func (timer *Timer) run() {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for range ticker.C {
		timer.mu.Lock()
		if timer.val > 0 {
			timer.val--
			if timer.action != nil {
				timer.action()
			}
		}
		timer.mu.Unlock()
	}
}

func (timer *Timer) Set(v byte) {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	timer.val = v
}

func (timer *Timer) Get() byte {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	return timer.val
}
