package perftest

import (
	"errors"
	"sync"
	"time"
)

// Manager manages workers and has a central store for test info
type Manager struct {
	s *store

	// Even mutex protected, mostly read lock
	sync.RWMutex
	workerMap map[string]*worker
}

// Create a new Manager
func Create() *Manager {
	tm := new(Manager)
	tm.s = new(store)
	tm.workerMap = make(map[string]*worker)

	return tm
}

// Add Adds a test
func (tm *Manager) Add(testID string, t Params) {
	w := createWorker(tm, t)
	go w.run()

	w.Request <- struct{}{}
	<-w.Response

	tm.Lock()
	defer tm.Unlock()
	tm.workerMap[testID] = w
}

// Get the test result using testID
func (tm *Manager) Get(testID string) (Result, error) {
	if r, e := tm.s.get(testID); e == nil {
		// if we have the result in the store, the next request thread is
		// responsible for closing the worker go-routine
		tm.RLock()
		defer tm.RUnlock()
		// the request thread never blocks and do it in best effort manner
		// so we have a race condition of mutliple request thread shutting
		// down the worker only one wins the others timeout
		if w, ok := tm.workerMap[testID]; ok {
			go func() {
				select {
				case w.Exit <- struct{}{}:
				case <-time.After(5 * time.Second):
				}
			}()
		}
		return r.Result, nil
	}

	tm.RLock()
	defer tm.RUnlock()
	if w, ok := tm.workerMap[testID]; ok {
		w.Request <- struct{}{}
		r := <-w.Response
		return r, nil
	}

	return nil, errors.New("test doesn't exist")
}
