package perftest

import (
	"log"
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
	// create a buffered channel so no matter what worker is doing, it will receive
	// the signal from Exit channel afterwards
	w.Exit = make(chan struct{}, 1)
	w.dbIDTracker = new(DBIDTracker)
	w.sc = t.Controller()
	if w.sc == nil {
		panic("ERR: Cannot create worker with nil controller")
	}

	if e := w.sc.UpdateBaselineIDs(w.dbIDTracker); e != nil {
		log.Fatalf("ERR: update failed, %v", e)
	}

	w.tm = tm

	var tinfo TestInfo
	tinfo.Params = t
	var tr TestResult
	tr.StartTime = time.Now()
	tr.Done = false
	tr.AdditionalInfo = t.Info()
	tr.Keywords = t.Keywords()
	tr.SetCPUMax(float64(0))
	tr.SetMemMax(float64(0))

	if e := w.sc.UpdateDBParameters(t.DBConfig().Database, &(tr.DBParam)); e != nil {
		log.Fatalf("ERR: update system parameters failed: %v", e)
	}

	switch t.(type) {
	default:
		panic("ERR: Unknown test parameter type while creating worker thread")

	case *RatingParams:
		w.tt = RATING

		rr := new(RatingResult)
		rr.FilesCompleted = 0
		rr.MinRate = 0
		rr.Rates = make([]float32, 0)
		rr.TestResult = tr
		tinfo.Result = rr

	case *BillingParams:
		w.tt = BILLING

		rr := new(BillingResult)
		rr.UserPackageBillRate = make([]uint32, 0)
		rr.TestResult = tr
		tinfo.Result = rr
	}

	w.ti = &tinfo
	return w
}

// if the test is completed, update the store too
func (w *worker) update() {
	if w.ti.Result.Result().Done {
		return
	}

	w.trackKPI()

	switch w.tt {
	case RATING:
		if e := w.sc.UpdateRatingResult(w.ti, w.dbIDTracker); e != nil {
			log.Fatalf("ERR: Worker failed updating rating results, %v", e)
		}

	case BILLING:
		if e := w.sc.UpdateBillingResult(w.ti, w.dbIDTracker); e != nil {
			log.Fatalf("ERR: Worker failed updating billing results, %v", e)
		}

	default:
	}

	return
}

//
func (w *worker) sendResult() {
	w.Response <- w.ti.Result
}

func (w *worker) trackKPI() {
	w.sc.TrackKPI(w.ti.Result)
}

func (w *worker) run() {
	wt := w.ti.Params.CollectionInterval()
	if wt != 0 {
		waitTime = wt
	}
	timer := time.NewTimer(waitTime)
	for {
		if w.ti.Result.Result().Done {
			w.once.Do(func() {
				w.tm.s.add(w.ti.Params.TestID(), w.ti)
				timer.Stop()
			})
		}

		select {
		case <-timer.C:
			w.update()
			timer.Reset(waitTime)
		case <-w.Request:
			// TODO not to update if we calculate the rates by counting the rows
			//w.update()
			w.sendResult()
			// never reset timer if we want calcuate by counting the rows
			// return false if timer is stopped
			//timer.Reset(waitTime)
		case <-w.Exit:
			w.tm.Lock()
			defer w.tm.Unlock()
			delete(w.tm.workerMap, w.ti.Params.TestID())
			close(w.Response)
			return
		}
	}
}
