package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/perf-prototype/stats"
)

// TestInfo stores all the test related information
type TestInfo struct {
	Params *Params
	Result interface{}
}

// store to save test information
type store struct {
	sync.RWMutex
	info map[string]*TestInfo
}

// TestManager stores the test information
type TestManager struct {
	db *stats.Controller
	s  *store
}

func (s *store) add(uuid string, t *TestInfo) {
	s.Lock()
	defer s.Unlock()
	s.info[uuid] = t
}

// get a copy of testInfo from the store
func (s *store) get(uuid string) (TestInfo, error) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.info[uuid]; !ok {
		return TestInfo{}, errors.New("Test doesn't exist.")

	}

	return *s.info[uuid], nil
}

func (s *store) update(uuid string, t *TestResult) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.info[uuid]; !ok {
		log.Println("WARNING: updating an non-existing test.")
		return errors.New("Update non-existing test.")
	}

	tm := *t
	s.info[uuid].testResult = &tm
	return nil
}

// RateResult to save the rate information
type RatingTestResult struct {
	*TestResult
	Rates          []float32 `json:"rates"`
	MinRate        float32   `json:"MIN_rate"`
	AvgRate        float32   `json:"AVG_rate"`
	FilesCompleted uint      `json:"number_of_files_completed"`
}

// TestResult to store all related results
type TestResult struct {
	StartTime      time.Time         `json:"-"`
	LastLog        uint64            `json:"-"`
	Done           bool              `json:"test_completed"`
	AdditionalInfo map[string]string `json:"additional_info"`
}

// Params to hold the testing parameters
type Params struct {
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
