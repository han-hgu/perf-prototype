package rating

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// variant for the amount field in the record generated
const delta uint32 = 10

// testInfo stores all the test related information
type testInfo struct {
	testParams TestParams
	testResult testResult
}

// store to save all test related information
// TODO: save result in disk, purge memory
type store struct {
	sync.RWMutex
	info map[string]*testInfo
}

func (s *store) add(uuid string, t *testInfo) {
	s.Lock()
	defer s.Unlock()
	s.info[uuid] = t
}

// get a copy of testInfo from the store
func (s *store) get(uuid string) (testInfo, error) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.info[uuid]; !ok {
		return testInfo{}, errors.New("Test doesn't exist")
	}

	return *s.info[uuid], nil
}

/*
controller responsibility
    - Spawn threads to do the record dropping in the desired location
    - UUID for book keeping
    - Log and publish stats
*/
type controller struct {
	// TODO: standalone book-keeping service, controller
	// is only responsible of delegating the tasks

	// s to store test information
	s *store
}

// testResult to store all related results
type testResult struct {
	startTime time.Time
}

// TestParams to hold the testing parameters
type TestParams struct {
	AmtFieldIndex       int      `json:"amount_field_index"`
	TimpstampFieldIndex int      `json:"timestamp_field_index"`
	RecordSizePerFile   int      `json:"records_per_file"`
	NumOfFiles          int      `json:"number_of_files"`
	NumRecordsPerFile   int      `json:"number_of_records_per_file"`
	RawFields           []string `json:"raw_fields"`
	UseExistingFile     bool     `json:"use_existing_file"`
	DropLocation        string   `json:"drop_location"`
	FilenamePrefix      string   `json:"filename_prefix"`
}

var c *controller
var once sync.Once

// StartProcess to start the rating test
// ownership of t is now transferred to the stats controller
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

	if t.UseExistingFile {

	} else {
		if t.FilenamePrefix == "" {
			t.FilenamePrefix = uid
		}
		if e := c.createFile(t); e != nil {
			return "", e
		}
	}

	// construct testinfo object
	tinfo := testInfo{
		testParams: *t,
		testResult: *new(testResult),
	}

	c.s.add(uid, &tinfo)
	return uid, nil
}

// exists returns true if file path exists
func exists(path string) error {
	// TODO: there are other errors besides the file doesn't exist error
	_, err := os.Stat(path)
	return err
}

func (rc *controller) createFile(t *TestParams) error {
	// check to see if the location exist, location specified must exist
	if err := exists(t.DropLocation); err != nil {
		return err
	}

	var filename string
	for i := 0; i < t.NumOfFiles; i++ {
		filename = t.DropLocation + "/" + t.FilenamePrefix + "-" + strconv.Itoa(i) + ".csv"

		fo, err := os.Create(filename)
		defer func() {
			if e := fo.Close(); e != nil {
				panic(e)
			}
		}()

		if err != nil {
			return err
		}

		for i := 0; i < t.NumRecordsPerFile; i++ {
			// No random, rate repeatly using the current timestamp for phase 1
			// 20060102150405 is const have to specify it this way, refer to
			// http://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
			tns := time.Now().Format("20060102150405.000")

			// replace the timestamp
			t.RawFields[t.TimpstampFieldIndex] = tns
			fo.WriteString(strings.Join(t.RawFields, ",") + "\n")
		}
	}

	return nil
}
