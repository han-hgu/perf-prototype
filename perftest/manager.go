package perftest

import (
	"fmt"
	"time"

	"github.com/perf-prototype/stats"
)

// Params interface to abstract out params
type Params interface {
	GetInfo() map[string]string
}

// Result interface to abstract out results
type Result interface {
	GetResult() *TestResult
}

// TestInfo stores all the test related information
type TestInfo struct {
	Params Params
	Result Result
	w      *worker
}

// RatingParams to hold the testing parameters
type RatingParams struct {
	AmtFieldIndex       int               `json:"amount_field_index"`
	TimpstampFieldIndex int               `json:"timestamp_field_index"`
	NumOfFiles          int               `json:"number_of_files"`
	NumRecordsPerFile   int               `json:"number_of_records_per_file"`
	RawFields           []string          `json:"raw_fields"`
	UseExistingFile     bool              `json:"use_existing_file"`
	DropLocation        string            `json:"drop_location"`
	FilenamePrefix      string            `json:"filename_prefix"`
	StartingEventlogID  int               `json:"starting_id"`
	AdditionalInfo      map[string]string `json:"additional_info"`
}

// GetInfo to integrate RatingParams to the Params interface
func (rp *RatingParams) GetInfo() map[string]string {
	return rp.AdditionalInfo
}

// TestResult to store generic results
type TestResult struct {
	StartTime      time.Time         `json:"-"`
	LastLog        uint64            `json:"-"`
	Done           bool              `json:"test_completed"`
	AdditionalInfo map[string]string `json:"additional_info"`
}

// RatingResult to save the rate information
type RatingResult struct {
	*TestResult
	Rates          []float64 `json:"rates"`
	MinRate        float32   `json:"MIN_rate"`
	AvgRate        float32   `json:"AVG_rate"`
	FilesCompleted int       `json:"number_of_files_completed"`
}

// GetResult to integrate RatingResult to the Result interface
func (rr *RatingResult) GetResult() *TestResult {
	return rr.TestResult
}

// Manager stores the test information
type Manager struct {
	db *stats.Controller
	s  *store
}

// Create a new Manager
func Create(dbc *stats.Controller) *Manager {
	tm := new(Manager)
	s := new(store)
	tm.s = s
	tm.s.info = make(map[string]*TestInfo)
	tm.db = dbc

	return tm
}

// Add a test to manager
func (tm *Manager) Add(testID string, t *TestInfo) {
	tm.s.add(testID, t)
	w := createWorker()
	w.run(testID, tm)
}

// Get the test result using testID
func (tm *Manager) Get(testID string) (Result, error) {
	var r Result

	ti, e := tm.s.get(testID)
	fmt.Println("HAN >>>>> ti.result in code", ti.Result, e)
	if e != nil {
		return nil, e
	}

	// ti.w can only be set by the manager so no race condition
	if ti.w != nil {
		ti.w.Request <- struct{}{}
		r = <-ti.w.TestResult
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
	} else { // worker is unregistered(results loaded from db) or deregistered
		r = ti.Result
	}

	return r, nil
}
