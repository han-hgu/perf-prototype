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
	switch ptype := t.(type) {
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

func (w *worker) update() (testComplete bool) {
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
	for {
		select {
		case <-timer.C:
			if testCompleted := w.update(); testCompleted {
				timer.Stop()
				w.tm.s.add(w.ti.Params.GetTestID(), w.ti)

				// Now that the test result is published to the
				// store, terminate
				return
			}
			timer.Reset(waitTime)
		case <-w.Request:
			if testCompleted := w.update(); testCompleted {
				timer.Stop()
				w.sendResult()
				w.tm.s.add(w.ti.Params.GetTestID(), w.ti)

				// Now that the test result is published to the
				// store, terminate
				return
			}
			w.sendResult()
			timer.Reset(waitTime)
		}
	}
}
