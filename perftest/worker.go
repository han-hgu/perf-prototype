package perftest

import (
	"log"
	"time"
)

// testType for test type
type testType uint8

// enum for test type
const (
	RATING testType = iota + 1
	BILLING
)

var waitTime = 10 * time.Second

// Worker for stats handling and background sync
type worker struct {
	tt          testType
	tm          *Manager
	ti          *TestInfo
	Request     chan struct{} // request for test info, ingress
	Response    chan Result   // response from worker for test info
	Exit        chan struct{}
	dbIDTracker *DBIDTracker
}

func createWorker(tm *Manager, t Params) *worker {
	w := new(worker)
	w.Request = make(chan struct{})
	w.Response = make(chan Result)
	w.Exit = make(chan struct{})
	w.dbIDTracker = new(DBIDTracker)
	if e := tm.db.UpdateBaselineIDs(w.dbIDTracker); e != nil {
		log.Fatalf("ERR: update failed, %v", e)
	}
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

// if the test is completed, update the store too
func (w *worker) update() {
	if w.ti.Result.GetResult().Done {
		w.tm.s.add(w.ti.Params.GetTestID(), w.ti)
		return
	}

	switch w.tt {
	case RATING:
		if e := w.tm.db.UpdateRatingResult(w.ti, w.dbIDTracker); e != nil {
			log.Fatalf("Worker error updating rating results, %v", e)
		}

	default:
	}

	return
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
			w.update()
			timer.Reset(waitTime)
		case <-w.Request:
			// choose not to update
			//w.update()
			w.sendResult()
			//timer.Reset(waitTime)
		case <-w.Exit:
			timer.Stop()
			return
		}
	}
}
