package perftest

import (
	"errors"
	"sync"
)

// Manager interacts with the stats collector and stores all
// test information
type Manager struct {
	db iController
	s  *store

	// Even mutex protected, mostly read lock
	sync.RWMutex
	workerMap map[string]*worker
}

// Create a new Manager
func Create(dbc iController) *Manager {
	tm := new(Manager)
	tm.s = new(store)
	tm.db = dbc
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
