package rating

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

// variant for the amount field in the record generated
const delta uint32 = 10

/*
ratingController is responsible of
    * spawn threads to do the record dropping in the desired lolcation
    * Assign UUID for each rating test for book keeping
    * log and publish real-time stats
*/

//testInfo stores all the
type testInfo struct {
	testParams *TestParams
	testResult *testResult
}

// store to save all test related info, for now it is in memory and no
// purging
type store struct {
	sync.RWMutex
	info map[string]*testInfo
}

func (s *store) add(uuid string, t *testInfo) {
	s.Lock()
	defer s.Unlock()
	s.info[uuid] = t
}

func (s *store) get(uuid string) *testResult {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.info[uuid]; !ok {
		return nil
	}
	return s.info[uuid].testResult
}

type controller struct {
	// stores existing test info
	// TODO: need a standalone bookkeeping service for this, controller
	// is only responsible of delegating the tasks

	s *store
}

// testResult to store all related results
type testResult struct {
	startTime time.Time
}

// TestParams to hold the testing parameters
type TestParams struct {
	AmtFieldIndex       uint32   `json:"amount_field_index"`
	TimpstampFieldIndex uint32   `json:"timestamp_field_index"`
	BatchSize           uint32   `json:"batch_size"`
	NumOfFiles          uint32   `json:"number_of_files"`
	RawFields           []string `json:"raw_fields"`
	DropLocation        string   `json:"drop_location"`
	FilenamePrefix      string   `json:"-"`
}

var c *controller
var once sync.Once

// StartProcess to start the rating test
func StartProcess(t *TestParams) (id string, err error) {
	once.Do(func() {
		c = &controller{}
		c.s = new(store)
		c.s.info = make(map[string]*testInfo)
	})

	// allocate uuid for the test run
	uid, e := newUUID()
	if e != nil {
		return "", errors.New("fail to generate UUID")
	}

	t.FilenamePrefix = uid
	if e := c.fileDrop(t); e != nil {
		return "", e
	}

	// construct testinfo object
	tinfo := &testInfo{
		testParams: t,
		testResult: new(testResult),
	}

	c.s.add(uid, tinfo)

	return uid, nil
}

// check if file path exists
func exists(path string) error {
	// TODO: there are other errors besides the file doesn't exist error
	_, err := os.Stat(path)
	return err
}

func (rc *controller) fileDrop(t *TestParams) error {
	// check to see if the location exist, location specified must exist
	if err := exists(t.DropLocation); err != nil {
		return err
	}

	// TODO: for Phase 1 use the UUID as the file name, NumOfFiles is always set to 1
	filename := t.DropLocation + "/" + t.FilenamePrefix + ".csv"
	fo, err := os.Create(filename)
	defer func() {
		if e := fo.Close(); e != nil {
			panic(e)
		}
	}()

	if err != nil {
		return err
	}

	for i := uint32(0); i < t.BatchSize; i++ {
		// No random, rate repeatly using the current timestamp for phase 1
		// 20060102150405 is const have to specify it this way, refer to
		// http://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
		tns := time.Now().Format("20060102150405.000")

		// replace the timestamp
		t.RawFields[t.TimpstampFieldIndex] = tns
		fo.WriteString(strings.Join(t.RawFields, ",") + "\n")
	}

	return nil
}
