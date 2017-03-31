package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/perf-prototype/perftest"
	"github.com/perf-prototype/stats"
)

type controller struct {
	tm *perftest.Manager
	sc *stats.Controller
}

var c *controller
var once sync.Once

// singleton
func initController() {
	once.Do(func() {
		c = &controller{}

		var conf stats.DBConfig
		if _, err := toml.DecodeFile("perf.conf", &conf); err != nil {
			log.Fatal(err)
		}

		sc := stats.CreateController(&conf)
		c.sc = sc
		c.tm = perftest.Create(sc)
	})
}

// GetResult returns the test result based on the UUID
func GetResult(testID string) (perftest.Result, error) {
	// query before any test is started
	initController()
	return c.tm.Get(testID)
}

// StartRateTest starts a rating test
func StartRateTest(t *perftest.RatingParams) (id string, err error) {
	initController()

	// allocate uuid for the test run
	uid, e := newUUID()
	if e != nil {
		return "", errors.New("fail to generate test ID")
	}

	if t.UseExistingFile {

	} else {
		if t.FilenamePrefix == "" {
			t.FilenamePrefix = uid
		}
		if e := createFile(t); e != nil {
			return "", e
		}
	}

	// make a copy, don't bother tracking ownership
	tp := *t
	tr := new(perftest.RatingResult)
	tr.TestResult = new(perftest.TestResult)

	tinfo := perftest.TestInfo{
		Params: &tp,
		Result: tr,
	}

	tr.Done = false
	tr.StartTime = time.Now()
	tr.FilesCompleted = 0

	if tp.StartingEventlogID == 0 {
		tr.LastLog = c.sc.GetLastIDFromEventLog("")
	} else {
		tr.LastLog = uint64(tp.StartingEventlogID)
	}

	tr.AdditionalInfo = tp.AdditionalInfo

	c.tm.Add(uid, &tinfo)
	return uid, nil
}
