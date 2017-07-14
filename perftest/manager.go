package perftest

import (
	"errors"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Manager to manage workers and the store for finished tests
type Manager struct {
	s storage

	// Even mutex protected, mostly read lock
	sync.RWMutex
	workerMap map[bson.ObjectId]*worker
}

// Create a new manager
func Create(s storage) *Manager {
	tm := new(Manager)
	tm.s = s
	tm.workerMap = make(map[bson.ObjectId]*worker)

	return tm
}

// Teardown to tear down the manager properly
func (tm *Manager) Teardown() {
	tm.s.Teardown()
}

// Add a test
func (tm *Manager) Add(testID bson.ObjectId, t Params) {
	w := createWorker(tm, t)
	go w.run()

	w.Request <- struct{}{}
	<-w.Response

	tm.Lock()
	defer tm.Unlock()
	tm.workerMap[testID] = w
}

// Get the test result providing testID
func (tm *Manager) Get(testID bson.ObjectId) (Result, error) {
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
		return r, nil
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

// return true if all slice elements of A can be found in B(A is a sub-slice of B)
func contains(A []string, B []string) bool {
	keys := make(map[string]struct{})
	for _, v := range B {
		keys[v] = struct{}{}
	}

	for _, v := range A {
		if _, ok := keys[v]; !ok {
			return false
		}
	}

	return true
}

// GetAll returns all test meta data with the provided tags, if tags is nil,
// return all test meta data
func (tm *Manager) GetAll(tags []string) ([]TestResultSV, error) {
	retVal := make([]TestResultSV, 0)

	tm.RLock()
	for _, w := range tm.workerMap {
		trsv := TestResultSV{}
		w.Request <- struct{}{}
		m := <-w.Response
		if contains(tags, m.MetaData().Keywords) {
			trsv.ID = m.TestID()
			trsv.Md = m.MetaData()
			retVal = append(retVal, trsv)
		}
	}
	tm.RUnlock()

	rv, err := tm.s.getTestResultSVByTags(tags)

	if err != nil {
		return nil, err
	}

	retVal = append(retVal, rv...)
	return retVal, nil
}
