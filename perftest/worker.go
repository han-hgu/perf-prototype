package perftest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

var waitTime = 30 * time.Second

// want the http client to not affect the worker at all so imminent timeout
const appClientReqTimeout = 300 * time.Millisecond

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
	appStatsC   *http.Client
}

func createWorker(tm *Manager, t Params) *worker {
	w := new(worker)
	w.Request = make(chan struct{})
	w.Response = make(chan Result)

	// initialize the http client for app server stats
	if w.appStatsC == nil {
		w.appStatsC = &http.Client{Timeout: appClientReqTimeout}
	}

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
	if t.CollectionInterval() == "" {
		tr.CInterval = waitTime.String()
	} else {
		tr.CInterval = t.CollectionInterval()
	}

	// copy parameters to results
	tr.StartTime = time.Now()
	tr.ID = t.TestID()
	tr.CTitle = t.ChartTitle()
	tr.Done = false
	tr.Cmt = t.Comment()
	tr.Keywords = t.Keywords()

	// Update app parameters
	tr.AppConf = *t.AppConfig()

	if e := w.sc.UpdateDBParameters(t.DBConfig(), &(tr.DBParams)); e != nil {
		log.Fatalf("ERR: update system parameters failed: %v", e)
	}

	switch t.(type) {
	default:
		panic("ERR: Unknown test type while creating worker thread")

	case *RatingParams:
		w.tt = RATING

		rr := new(RatingResult)
		rr.FilesCompleted = 0
		rr.MinRate = 0
		rr.Rates = make([]float32, 0)
		rr.TestResult = tr
		tinfo.Result = rr
		rr.Type = "rating"

	case *BillingParams:
		w.tt = BILLING

		rr := new(BillingResult)
		rr.OwnerName = t.(*BillingParams).OwnerName
		rr.UserPackageBillRate = make([]uint32, 0)
		rr.TestResult = tr
		tinfo.Result = rr
		rr.Type = "billing"
	}

	w.ti = &tinfo

	// baseline the database KPI
	// TODO find a better place and multitasking
	cpudontcare := new(float32)
	lreadbase := &(w.ti.Result.DBServerStats().LReadsBase)
	lwritebase := &(w.ti.Result.DBServerStats().LWritesBase)
	preadbase := &(w.ti.Result.DBServerStats().PReadsBase)
	w.sc.TrackKPI(nil, tinfo.Params.DBConfig().Database, cpudontcare, lreadbase, lwritebase, preadbase)

	return w
}

func (w *worker) TrackKPI() {
	var (
		appCPU   = new(float32)
		appMem   = new(float32)
		dbsysCPU = new(float32)
		dbCPU    = new(float32)
		lreads   = new(uint64)
		lwrites  = new(uint64)
		preads   = new(uint64)
	)

	kpiStart := time.Now()

	var wg sync.WaitGroup
	// Track app server KPI
	wg.Add(1)
	go w.TrackAppServerKPI(&wg, appCPU, appMem)

	// Track database server KPI
	wg.Add(1)
	go w.sc.TrackKPI(&wg, w.ti.Params.DBConfig().Database, dbCPU, lreads, lwrites, preads)

	wg.Add(1)
	go w.TrackDBSysCPU(&wg, dbsysCPU)

	wg.Wait()
	log.Printf("INFO: KPI tracking takes %v to complete.\n", time.Since(kpiStart))

	if *lreads >= w.ti.Result.DBServerStats().LReadsBase {
		w.ti.Result.AddLogicalReads(*lreads - w.ti.Result.DBServerStats().LReadsBase)
	} else {
		// Reset the base since the system resets the stats in dm_exec_query_stats
		w.ti.Result.AddLogicalReads(0)
	}
	w.ti.Result.DBServerStats().LReadsBase = *lreads

	if *lwrites >= w.ti.Result.DBServerStats().LWritesBase {
		w.ti.Result.AddLogicalWrites(*lwrites - w.ti.Result.DBServerStats().LWritesBase)
	} else {
		w.ti.Result.AddLogicalWrites(0)
	}
	w.ti.Result.DBServerStats().LWritesBase = *lwrites

	if *preads >= w.ti.Result.DBServerStats().PReadsBase {
		w.ti.Result.AddPhysicalReads(*preads - w.ti.Result.DBServerStats().PReadsBase)
	} else {
		w.ti.Result.AddPhysicalReads(0)
	}
	w.ti.Result.DBServerStats().PReadsBase = *preads

	w.ti.Result.AddDBCPU((*dbsysCPU) * (*dbCPU) / 100)
	w.ti.Result.AddAppServerCPU(*appCPU)
	w.ti.Result.AddAppServerMem(*appMem)
}

// TrackAppServerKPI to track the app server stats, minimum impact to the app
// server even the app server micro service is down
func (w *worker) TrackAppServerKPI(wg *sync.WaitGroup, cpu *float32, mem *float32) {
	if wg != nil {
		defer wg.Done()
	}

	var pfstats PerfMonStats
	rsp, e := w.appStatsC.Get(w.ti.Params.AppConfig().URL)
	if e != nil {
		fmt.Printf("WARNING: failed to get app server stats from %v, error: %v\n", w.ti.Params.AppConfig().URL, e)
	} else {
		json.NewDecoder(rsp.Body).Decode(&pfstats)
	}

	*cpu = pfstats.CPU
	*mem = pfstats.Mem
}

// TrackDBSysCPU takes the CPU usage of the database server, this is the total
// CPU used by the server
func (w *worker) TrackDBSysCPU(wg *sync.WaitGroup, cpu *float32) {
	if wg != nil {
		defer wg.Done()
	}

	var pfstats PerfMonStats
	rsp, e := w.appStatsC.Get(w.ti.Params.DBConfig().URL)
	if e != nil {
		fmt.Printf("WARNING: failed to get database stats from %v, error: %v\n", w.ti.Params.AppConfig().URL, e)
	} else {
		json.NewDecoder(rsp.Body).Decode(&pfstats)
	}

	*cpu = pfstats.CPU
}

func (w *worker) update() {
	if w.ti.Result.Result().Done {
		return
	}

	w.TrackKPI()

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

func (w *worker) sendResult() {
	w.Response <- w.ti.Result
}

func (w *worker) run() {
	wt := w.ti.Params.CollectionInterval()
	if wt != "" {
		var err error
		waitTime, err = time.ParseDuration(wt)
		if err != nil {
			log.Fatalf("Invalid duration %s", wt)
		}
	}

	timer := time.NewTimer(waitTime)
	for {
		if w.ti.Result.Result().Done {
			w.once.Do(func() {
				w.tm.s.add(w.ti.Result)
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
