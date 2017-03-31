package perftest

import (
	"log"
	"time"
)

const waitTime time.Duration = 1 * time.Second

// Worker for stats handling and background sync
type worker struct {
	TestID     string
	Request    chan struct{}
	Exit       chan struct{}
	TestResult chan Result
}

func createWorker() *worker {
	w := new(worker)
	w.Request = make(chan struct{})
	w.Exit = make(chan struct{})
	w.TestResult = make(chan Result)
	return w
}

func (w *worker) update(testID string, tm *Manager) {
	ti, err := tm.s.get(testID)
	if err != nil {
		log.Println("ERROR: worker update failed getting the testInfo")
		return
	}

	switch tp := ti.Params.(type) {
	default:
		return
	case *RatingParams:
		// TODO: panic if cast fails
		tr := ti.Result.(*RatingResult)
		if tr.FilesCompleted == tp.NumOfFiles {
			return
		}
		// always use the store to do the update
		// make a copy of testInfo to point to this new result
		trn := *tr
		upid, fprocessed, r := tm.db.GetRates(tp.FilenamePrefix, tr.LastLog)
		trn.FilesCompleted += fprocessed
		trn.LastLog = upid
		trn.Rates = append(trn.Rates, r...)

		// set the Done flag so that the manager can deregister the worker
		if trn.FilesCompleted >= tp.NumOfFiles {
			trn.Done = true
		}
		ti.Result = &trn
		tm.s.update(testID, ti)
	}
}

//
func (w *worker) sendResult(testID string, tm *Manager) {
	ti, e := tm.s.get(testID)
	if e != nil {
		panic(e)
	}

	// This should always be place at last to solve the race condition that the
	// server receives a request while the worker is shutting down
	w.TestResult <- ti.Result
}

func (w *worker) sync(testID string, tm *Manager) {
	timer := time.NewTimer(waitTime)
	for {
		select {
		case <-w.Exit:
			w.sendResult(testID, tm)
			return
		case <-timer.C:
			w.update(testID, tm)
			timer.Reset(waitTime)
		case <-w.Request:
			w.update(testID, tm)
			w.sendResult(testID, tm)
			timer.Reset(waitTime)
		}
	}
}

// Run starts the sync process
func (w *worker) run(testID string, tm *Manager) error {
	ti, err := tm.s.get(testID)
	if err != nil {
		log.Fatal(err)
	}

	go w.sync(testID, tm)

	// send a request to test the backend sync process, register the worker
	// only if the goroutine is in the for loop, it is possible that a request
	// is landed after the worker is created but before the goroutine is about
	// to enter the for loop, no one is able to answer the get request in this
	// case
	w.Request <- struct{}{}
	<-w.TestResult
	ti.w = w
	if e := tm.s.update(testID, ti); e != nil {
		return e
	}
	return nil
}
