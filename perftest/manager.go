package perftest

import (
	"errors"
	"sync"
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
