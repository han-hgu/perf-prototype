package perftest

import (
	"time"
)

// testType for test type
type testType uint8

// enum for test type
const (
	RATING testType = iota + 1
	BILLING
)

const waitTime time.Duration = 1 * time.Second

// Worker for stats handling and background sync
type worker struct {
	tt       testType
	tm       *Manager
	ti       *TestInfo
	Request  chan struct{} // request for test info, ingress
	Response chan Result   // response from worker for test info
}

func createWorker(tm *Manager, t Params) *worker {
	w := new(worker)
	w.Request = make(chan struct{})
	w.Response = make(chan Result)
	w.tm = tm

	var tinfo TestInfo
	// create a testInfo obj
	switch t.(type) {
	default:
		panic("ERR: Unknown test parameter type while creating worker thread")

	case *RatingParams:
		w.tt = RATING

		rr := new(RatingResult)
		rr.Done = false
		rr.StartTime = time.Now()
		rr.FilesCompleted = 0
		rr.AdditionalInfo = t.GetInfo()

		tinfo.Params = t
		tinfo.Result = rr
	}

	w.ti = &tinfo
	return w
}

func (w *worker) update() (testCompleted bool) {
	return false
	/*
		switch tp := ti.Params.(type) {
		default:
			return
		case *RatingParams:
			// TODO: panic if cast fails
			tr, ok := ti.Result.(*RatingResult)

			if tr.Done {
				return
			}

			// always use the store to do the update
			// make a copy of testInfo to point to this new result
			trn := *tr
			upid, fprocessed, r := tm.db.GetRates(tp.FilenamePrefix, tr.LastLog)
			trn.FilesCompleted += fprocessed
			trn.LastLog = upid
			trn.Rates = append(trn.Rates, r...)

			// set the Done flag so that the manager can de-register the worker
			if trn.FilesCompleted >= tp.NumOfFiles {
				trn.Done = true
			}
			ti.Result = &trn
			tm.s.update(testID, &ti)
		}
	*/
}

//
func (w *worker) sendResult() {
	w.Response <- w.ti.Result
}

func (w *worker) run() {
	timer := time.NewTimer(waitTime)
	testID := w.ti.Params.GetTestID()
	for {
		select {
		case <-timer.C:
			if testCompleted := w.update(); testCompleted {
				timer.Stop()

				// if there are goroutines visiting this testID, do nothing,
				// worker can't be terminated at this state
				nvp, ok := w.tm.visitorTracker[testID]
				if !ok {
					panic("tm.Add() should be called before run()")
				}

				// attempt to shutdown if 0 visitor, check-lock-check pattern
				if nvp.get() == 0 {
					w.tm.workerPoolLock.Lock()
					if nvp.get() == 0 {
						// it is possible the worker is not registered with
						// the workerPool yet and this will be a no-op
						delete(w.tm.workerPool, testID)
						delete(w.tm.visitorTracker, testID)
						w.tm.s.add(testID, w.ti)
						w.tm.workerPoolLock.Unlock()
						return
					}
					w.tm.workerPoolLock.Unlock()
				}
			}
			timer.Reset(waitTime)
		case <-w.Request:
			// design decision to only do the worker cleanup in the above timer
			// firing logic
			w.update()
			w.sendResult()
			timer.Reset(waitTime)
		}
	}
}
