package main

import (
	"sync"

	"gopkg.in/mgo.v2/bson"

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
		s := new(perftest.Store)
		s.Initialize()
		c.tm = perftest.Create(s)
	})
}

// Teardown to tear down the controller
func Teardown() {
	c.tm.Teardown()
}

// Result returns the test result based on the test id
func Result(testID string) (perftest.Result, error) {
	// query before any test is started
	initController()
	return c.tm.Get(bson.ObjectIdHex(testID))
}

// MetaData returns all test meta data
func TestResultSVs(tags []string) ([]perftest.TestResultSV, error) {
	initController()

	r, e := c.tm.GetAll(tags)
	return r, e
}

// StartRatingTest starts a rating test
func StartRatingTest(t *perftest.RatingParams) (id string, err error) {
	initController()

	// allocate uuid for the test run
	t.ID = bson.NewObjectId()

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
			t.FilenamePrefix = t.ID.String()
		}
		if e := createFile(t); e != nil {
			return "", e
		}
	}

	c.tm.Add(t.ID, t)
	return t.ID.Hex(), nil
}

// StartBillingTest starts a billing test
func StartBillingTest(t *perftest.BillingParams) (id string, err error) {
	initController()

	// allocate uuid for the test run
	t.ID = bson.NewObjectId()

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

	c.tm.Add(t.ID, t)
	return t.ID.String(), nil
}
