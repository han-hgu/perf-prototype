package perftest

import (
	"errors"
	"sync"
	"time"
)

// Manager to manage workers and the store for finished tests
type Manager struct {
	s *store

	// Even mutex protected, mostly read lock
	sync.RWMutex
	workerMap map[string]*worker
}

// Create a new manager
func Create() *Manager {
	tm := new(Manager)
	tm.s = new(store)
	tm.workerMap = make(map[string]*worker)

	return tm
}

// Add a test
func (tm *Manager) Add(testID string, t Params) {
	w := createWorker(tm, t)
	go w.run()

	w.Request <- struct{}{}
	<-w.Response

	tm.Lock()
	defer tm.Unlock()
	tm.workerMap[testID] = w
}

// Get the test result providing testID
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
				// This is a buffered channel, have to take into account that
				// the the worker is not able to handle the Exit within 5 sec
				case w.Exit <- struct{}{}:
				// If mutliple goroutines send to Exit channel, some of them will
				// block, this is to prevent the resource leak
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
