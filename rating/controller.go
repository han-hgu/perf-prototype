package rating

import (
	"errors"
	"sync"
	"time"

	"github.com/perf-prototype/stats"
)

// variant for the amount field in the record generated
const delta uint32 = 10

// testInfo stores all the test related information
type TestInfo struct {
	testParams *TestParams
	testResult *TestResult
}

// internal store saving test information
type store struct {
	sync.RWMutex
	info map[string]*TestInfo
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
		return TestInfo{}, errors.New("Test doesn't exist")
	}

	return *s.info[uuid], nil
}

/*
controller responsibility
    - Spawn threads for record dropping in the desired location
    - UUID for book keeping
    - Log and publish stats
*/
type controller struct {
	// TODO: standalone book-keeping service, controller
	// is only responsible of delegating the tasks

	// s to store test information
	s *store
}

// TestResult to store all related results
type TestResult struct {
	StartTime      time.Time `json:"-"`
	Done           bool      `json:"test_completed"`
	Rates          []float32 `json:"rates"`
	FilesCompleted uint      `json:"number_of_files_completed"`
	CurrLogID      uint64    `json:"-"`
}

// TestParams to hold the testing parameters
type TestParams struct {
	AmtFieldIndex       int      `json:"amount_field_index"`
	TimpstampFieldIndex int      `json:"timestamp_field_index"`
	NumOfFiles          int      `json:"number_of_files"`
	NumRecordsPerFile   int      `json:"number_of_records_per_file"`
	RawFields           []string `json:"raw_fields"`
	UseExistingFile     bool     `json:"use_existing_file"`
	DropLocation        string   `json:"drop_location"`
	FilenamePrefix      string   `json:"filename_prefix"`
}

var c *controller
var once sync.Once

func initController() {
	once.Do(func() {
		c = &controller{}
		c.s = new(store)
		c.s.info = make(map[string]*TestInfo)
	})
}

// QueryStats returns the test result from the test UUID
func QueryStats(testId string) *TestResult {
	// handle the case that the user could submit a GET request before
	// any actual task has started
	initController()
	stats.GetController()

	// find the testInfo from the store
	/* HAN >>>>
	t, found := c.s.info[testId]
	if !found {
		return nil
	}
	*/

	// find the id of the first record
	//dbc.GetIdFromEventLog("")
	return nil

	//return &(t.testResult)
}

// StartTest to start the rating test
func StartTest(t *TestParams) (id string, err error) {
	initController()

	// allocate uuid for the test run
	uid, e := newUUID()
	if e != nil {
		return "", errors.New("fail to generate UUID")
	}

	// get the latest eventlog ID before starting the test
	stats.GetController()

	if t.UseExistingFile {

	} else {
		if t.FilenamePrefix == "" {
			t.FilenamePrefix = uid
		}
		if e := createFile(t); e != nil {
			return "", e
		}
	}

	// Make a copy of testParam
	tnew := *t
	tinfo := TestInfo{
		testParams: &tnew,
		testResult: new(TestResult),
	}

	tinfo.testResult.Done = false
	tinfo.testResult.StartTime = time.Now()

	c.s.add(uid, &tinfo)
	return uid, nil
}
