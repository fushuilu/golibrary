package golibrary

import "sync"

type Async struct {
	mu *sync.RWMutex
}

func NewAsync() Async {
	return Async{
		mu: new(sync.RWMutex), // safety
	}
}

func (a *Async) Write(run func() error) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return run()
}

func (a *Async) Read(run func() error) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return run()
}
