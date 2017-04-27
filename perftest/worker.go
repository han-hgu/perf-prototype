package perftest

import (
	"log"
	"math"
	"sync"
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
	once        sync.Once
	sc          iController
}

func createWorker(tm *Manager, t Params) *worker {
	w := new(worker)
	w.Request = make(chan struct{})
	w.Response = make(chan Result)
	w.Exit = make(chan struct{})
	w.dbIDTracker = new(DBIDTracker)
	w.sc = t.GetController()
	if w.sc == nil {
		panic("ERR: Cannot create worker with nil controller")
	}

	if e := w.sc.UpdateBaselineIDs(w.dbIDTracker); e != nil {
		log.Fatalf("ERR: update failed, %v", e)
	}
	w.tm = tm

	// controller knows what actual db controller to use and should create the
	// instance already, all worker knows is the interface

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
		rr.MinRate = math.MaxFloat32

		tinfo.Params = t
		tinfo.Result = rr
	}

	w.ti = &tinfo
	return w
}

// if the test is completed, update the store too
func (w *worker) update() {
	if w.ti.Result.GetResult().Done {
		return
	}

	switch w.tt {
	case RATING:
		if e := w.sc.UpdateRatingResult(w.ti, w.dbIDTracker); e != nil {
			log.Fatalf("ERR: Worker failed updating rating results, %v", e)
		}

	default:
	}

	return
}

//
func (w *worker) sendResult() {
	w.Response <- w.ti.Result
}

func (w *worker) run() {
	timer := time.NewTimer(waitTime)
	for {
		if w.ti.Result.GetResult().Done {
			w.once.Do(func() {
				w.tm.s.add(w.ti.Params.GetTestID(), w.ti)
				timer.Stop()
			})
		}

		select {
		case <-timer.C:
			w.update()
			timer.Reset(waitTime)
		case <-w.Request:
			// TODO choose not to update?
			w.update()
			w.sendResult()
			// return false if timer is stopped
			timer.Reset(waitTime)
		case <-w.Exit:
			return
		}
	}
}
