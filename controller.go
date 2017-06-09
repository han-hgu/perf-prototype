package main

import (
	"errors"
	"sync"

	"github.com/perf-prototype/perftest"
	"github.com/perf-prototype/stats"
)

type controller struct {
	tm *perftest.Manager
}

var c *controller
var once sync.Once

// singleton
func initController() {
	once.Do(func() {
		c = &controller{}
		c.tm = perftest.Create()
	})
}

// Teardown to tear down the controller
func Teardown() {
	c.tm.Teardown()
}

// Result returns the test result based on the UUID
func Result(testID string) (perftest.Result, error) {
	// query before any test is started
	initController()
	return c.tm.Get(testID)
}

// MetaData returns all test meta data
func MetaData() []map[string]interface{} {
	initController()

	return c.tm.GetAll()
}

// StartRatingTest starts a rating test
func StartRatingTest(t *perftest.RatingParams) (id string, err error) {
	initController()

	// allocate uuid for the test run
	uid, e := newUUID()
	if e != nil {
		return "", errors.New("fail to generate test ID")
	}
	t.ID = uid

	// for rating test, controller creates and assigns the stats controller to t;
	// Perftest package should be flexible and only deal
	// with iController interface for future extensibility
	var statsDBConf stats.DBConfig
	statsDBConf.Server = t.DBConf.Server
	statsDBConf.Port = t.DBConf.Port
	statsDBConf.Database = t.DBConf.Database
	statsDBConf.UID = t.DBConf.UID
	statsDBConf.Pwd = t.DBConf.Pwd

	sc, err := stats.CreateController(&statsDBConf)
	if err != nil {
		return "", err
	}
	t.TestParams.DbController = sc

	if !t.UseExistingFile {
		if t.FilenamePrefix == "" {
			t.FilenamePrefix = uid
		}
		if e := createFile(t); e != nil {
			return "", e
		}
	}

	c.tm.Add(uid, t)
	return uid, nil
}

// StartBillingTest starts a billing test
func StartBillingTest(t *perftest.BillingParams) (id string, err error) {
	initController()

	// allocate uuid for the test run
	uid, e := newUUID()
	if e != nil {
		return "", errors.New("fail to generate test ID")
	}
	t.ID = uid

	// for rating test, controller creates and assigns the stats controller to t;
	// Perftest package should be flexible and only deal
	// with iController interface for future extensibility
	var statsDBConf stats.DBConfig
	statsDBConf.Server = t.DBConf.Server
	statsDBConf.Port = t.DBConf.Port
	statsDBConf.Database = t.DBConf.Database
	statsDBConf.UID = t.DBConf.UID
	statsDBConf.Pwd = t.DBConf.Pwd

	sc, err := stats.CreateController(&statsDBConf)
	if err != nil {
		return "", err
	}
	t.TestParams.DbController = sc

	c.tm.Add(uid, t)
	return uid, nil
}
