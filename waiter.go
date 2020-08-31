package kafkatesting

import (
	"fmt"
	"time"
)

const DefaultTimeout = 15 * time.Second

type Waiter struct {
	Timeout time.Duration
	Sleep   time.Duration
}

func NewDefaultWaiter() *Waiter {
	return &Waiter{
		Timeout: DefaultTimeout,
		Sleep:   100 * time.Millisecond,
	}
}

func (w Waiter) WaitFor(check func() (interface{}, error)) (interface{}, error) {
	deadline := time.Now().Add(w.Timeout)
	for {
		ret, err := check()

		if err == nil {
			return ret, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("Timeout reached after %d s: %w", w.Timeout/time.Second, err)
		}

		time.Sleep(w.Sleep)
	}
}