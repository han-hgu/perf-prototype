package perftest

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/perf-prototype/stats"
)

// Params interface to abstract out params
type Params interface {
	GetInfo() map[string]string
	GetTestID() string
}

// Result interface to abstract out results
type Result interface {
	GetResult() *TestResult
}

// TestInfo stores all the test related information
type TestInfo struct {
	Params Params
	Result Result
}

// TestParams to hold common test parameters for all test types
type TestParams struct {
	TestID         string            `json:"test_id"`
	AdditionalInfo map[string]string `json:"additional_info"`
}

// RatingParams to hold the testing parameters
type RatingParams struct {
	TestParams
	AmtFieldIndex          int           `json:"amount_field_index"`
	TimpstampFieldIndex    int           `json:"timestamp_field_index"`
	NumOfFiles             int           `json:"number_of_files"`
	NumRecordsPerFile      int           `json:"number_of_records_per_file"`
	RawFields              []string      `json:"raw_fields"`
	UseExistingFile        bool          `json:"use_existing_file"`
	DropLocation           string        `json:"drop_location"`
	FilenamePrefix         string        `json:"filename_prefix"`
	DataCollectionInterval time.Duration `json:"data_collection_interval"`
}

// GetInfo to integrate RatingParams to Params interface
func (rp *RatingParams) GetInfo() map[string]string {
	return rp.AdditionalInfo
}

// GetTestID returns the test ID
func (rp *RatingParams) GetTestID() string {
	return rp.TestID
}

// TestResult to store generic results
type TestResult struct {
	StartTime      time.Time         `json:"-"`
	LastEventLog   uint64            `json:"-"`
	Done           bool              `json:"test_completed"`
	AdditionalInfo map[string]string `json:"additional_info"`
}

// RatingResult to save the rate information
type RatingResult struct {
	TestResult
	Rates                 []float64 `json:"rates"`
	MinRate               float32   `json:"MIN_rate"`
	AvgRate               float32   `json:"AVG_rate"`
	UDRProcessed          uint64    `json:"udr_created"`
	UDRExceptionProcessed uint64    `json:"udr_exception_created"`
	FilesCompleted        int       `json:"files_completed"`
	LastUDRLog            uint64    `json:"-"`
	LastUDRExceptionLog   uint64    `json:"-"`
}

// GetResult to integrate RatingResult to Result interface
func (rr *RatingResult) GetResult() *TestResult {
	return &rr.TestResult
}

type atomicCounter int64

func (a *atomicCounter) increment() {
	atomic.AddInt64((*int64)(a), 1)
}

func (a *atomicCounter) decrement() {
	atomic.AddInt64((*int64)(a), -1)
}

func (a *atomicCounter) get() int64 {
	return atomic.LoadInt64((*int64)(a))
}

type testTracker struct {
	w *worker
	atomicCounter
}

// Manager interacts with the stats collector and stores all
// test information
type Manager struct {
	db *stats.Controller
	s  *store

	sync.RWMutex
	workerPool map[string]testTracker
}

// Create a new Manager
func Create(dbc *stats.Controller) *Manager {
	tm := new(Manager)
	tm.s = new(store)
	tm.db = dbc

	return tm
}

// Add creates a worker and tranfer the ownership of the
// testInfo to the worker, the worker is responsible until
// the test is completed
func (tm *Manager) Add(testID string, t Params) {
	w := createWorker(tm, t)
	go w.run()
	w.Request <- struct{}{}
}

// Get the test result using testID
func (tm *Manager) Get(testID string) (Result, error) {
	var r Result

	ti, e := tm.s.get(testID)
	if e != nil {
		return nil, e
	}

	if ti.w != nil {
		ti.w.Request <- struct{}{}
		r = <-ti.w.TestResult
		/*
			if r.GetResult().Done {
				// shutting down the worker syncing process
				ti.w.Exit <- struct{}{}
				// Get the last result, don't care about the value just for
				// synchronization purpose since worker will send the final result
				// and be blocked by the unbuffered channel, don't make the channel
				// unbuffered since it is possible the worker is doing cleanup and
				// server receives another
				<-ti.w.TestResult
				// set the worker to null, no race condition with the check in line
				// 99 since they are always the same thread
				ti.w = nil
			}
		*/
	} else { // worker is not registered(results loaded from db) or deregistered
		r = ti.Result
	}

	return r, nil
}
