package perftest

import (
	"testing"

	"github.com/go-test/deep"
	"gopkg.in/mgo.v2/bson"
)

func setupStoreTest() *Store {
	DBName = "test"
	CollectionName = "TestCollection"
	s := new(Store)
	s.Initialize()
	return s
}

func TestInitialize(t *testing.T) {
	s := setupStoreTest()
	defer s.Teardown()

	if err := s.mgoSession.Ping(); err != nil {
		t.Errorf("Store connection estabslished, %v\n", err)
	}
}

func TestStoreAdd(t *testing.T) {
	s := setupStoreTest()
	defer s.Teardown()

	tr := TestResult{}
	tr.Metadata.CInterval = "15s"
	tr.Metadata.CTitle = "my chart"
	tr.Keywords = []string{"abc"}
	tr.AddAppServerMem(1.3)
	tr.AddAppServerCPU(1.2)
	tr.AddDBServerCPU(1.1)
	tr.AddDBServerMem(100)
	tr.AddLogicalReads(1)
	tr.AddLogicalWrites(2)
	tr.AddPhysicalReads(3)
	tr.Metadata.AppConf.SysInfo = make(map[string]interface{}, 0)
	tr.Metadata.DBParams.SysInfo = make(map[string]interface{}, 0)
	tr.AppConf.Options = make(map[string]string)
	tr.ID = bson.NewObjectId()
	e := s.add(&tr)

	if e != nil {
		t.Errorf("Get an error while adding result to store, %v", e)
	}

	tra := TestResult{}
	s.mgoSession.DB(DBName).C(CollectionName).FindId(tr.ID).One(&tra)

	if diff := deep.Equal(tra, tr); diff != nil {
		t.Errorf("Result saved is the same as in memory, %v\n", diff)
	}

	// add the same record again
	if e := s.add(&tr); e != nil {
		t.Errorf("No error is returned if adding a duplicate: %v", e)
	}
	var results []TestResult
	s.mgoSession.DB(DBName).C(CollectionName).FindId(tr.ID).All(&results)

	if len(results) != 1 {
		t.Error("Only one record found after adding it multiple times to the store")
	}

	s.mgoSession.DB(DBName).C(CollectionName).RemoveAll(nil)
}

func TestStoreGet(t *testing.T) {
	s := setupStoreTest()
	defer s.Teardown()

	tr1 := TestResult{}
	tr1.Metadata.CInterval = "15s"
	tr1.Metadata.CTitle = "my chart"
	tr1.Metadata.Keywords = []string{"abc"}
	tr1.AddAppServerMem(1.3)
	tr1.AddAppServerCPU(1.2)
	tr1.AddDBServerCPU(1.1)
	tr1.AddDBServerMem(100)
	tr1.AddLogicalReads(1)
	tr1.AddLogicalWrites(2)
	tr1.AddPhysicalReads(3)
	tr1.AppConf.Options = make(map[string]string)
	tr1.ID = bson.NewObjectId()

	tr2 := TestResult{}
	tr2.Metadata.CInterval = "15s"
	tr2.Metadata.CTitle = "my chart"
	tr2.Metadata.Keywords = []string{"abc"}
	tr2.Metadata.AppConf.SysInfo = make(map[string]interface{}, 0)
	tr2.Metadata.DBParams.SysInfo = make(map[string]interface{}, 0)
	tr2.AddAppServerMem(1.3)
	tr2.AddAppServerCPU(1.2)
	tr2.AddDBServerCPU(1.1)
	tr2.AddDBServerMem(100)
	tr2.AddLogicalReads(1)
	tr2.AddLogicalWrites(2)
	tr2.AddPhysicalReads(3)
	tr2.AppConf.Options = make(map[string]string, 0)
	tr2.ID = bson.NewObjectId()
	s.add(&tr1)
	s.add(&tr2)

	// get a non-existent testID
	if _, err := s.get(bson.NewObjectId()); err == nil {
		t.Error("Get an error if a non-existenct testID is provided for get")
	}

	// get existing
	tra, err := s.get(tr2.ID)
	if err != nil {
		t.Error("Get no error if existing testID is provided for get")
	}

	if diff := deep.Equal(tra, &tr2); diff != nil {
		t.Errorf("Results stored is the same as provided, %v\n", diff)
	}

	s.mgoSession.DB(DBName).C(CollectionName).RemoveId(tr1.ID)
	s.mgoSession.DB(DBName).C(CollectionName).RemoveId(tr2.ID)
}

func TestGetTestResultSVByTags(t *testing.T) {
	s := setupStoreTest()
	defer s.Teardown()

	tr1 := RatingResult{}
	tr1.Metadata.CInterval = "15s"
	tr1.Metadata.CTitle = "my chart"
	tr1.Metadata.Keywords = []string{"abc"}
	tr1.AddAppServerMem(1.3)
	tr1.AddAppServerCPU(1.2)
	tr1.AddDBServerCPU(1.1)
	tr1.AddDBServerMem(100)
	tr1.AddLogicalReads(1)
	tr1.AddLogicalWrites(2)
	tr1.AddPhysicalReads(3)
	tr1.Type = "rating"
	tr1.UDRProcessed = 100
	tr1.AvgRate = 3.3
	tr1.Rates = []float32{1.0, 1.2, 1.3}
	tr1.ID = bson.NewObjectId()

	tr2 := TestResult{}
	tr2.Metadata.CInterval = "15s"
	tr2.Metadata.CTitle = "my chart2"
	tr2.Metadata.Keywords = []string{"def", "abc"}
	tr2.AddAppServerMem(1.3)
	tr2.AddAppServerCPU(1.2)
	tr2.AddDBServerCPU(1.1)
	tr2.AddDBServerMem(100)
	tr2.AddLogicalReads(1)
	tr2.AddLogicalWrites(2)
	tr2.AddPhysicalReads(3)
	tr2.ID = bson.NewObjectId()

	rr3 := RatingResult{}
	rr3.TestResult.Metadata.Keywords = []string{"rating"}
	rr3.ID = bson.NewObjectId()
	rr3.AppConf.Options = make(map[string]string)
	rr3.AppConf.Options["someOption1"] = "v1"
	rr3.AppConf.Options["someOption2"] = "v2"
	rr3.Metadata.AppConf.SysInfo = make(map[string]interface{}, 0)
	rr3.Metadata.DBParams.SysInfo = make(map[string]interface{}, 0)
	rr3SV := TestResultSV{}
	rr3SV.ID = rr3.ID
	rr3SV.Md = rr3.MetaData()

	br4 := BillingResult{}
	br4.TestResult.Metadata.Keywords = []string{"billing"}
	br4.ID = bson.NewObjectId()

	s.add(&tr1)
	s.add(&tr2)
	s.add(&rr3)
	s.add(&br4)
	tags := []string{"abc"}
	ret, e := s.getTestResultSVByTags(tags)
	if len(ret) != 2 || e != nil {
		t.Error("Correct number of results returned by tag")
	}

	tags = []string{}
	ret, e = s.getTestResultSVByTags(tags)
	if len(ret) != 4 || e != nil {
		t.Error("If tag is not provided, all test results are returned")
	}

	tags = []string{"rating"}
	ret, e = s.getTestResultSVByTags(tags)
	if len(ret) != 1 || e != nil {
		t.Error("Correct number of results returned by tag")
	}

	if diff := deep.Equal(rr3SV, ret[0]); diff != nil {
		t.Errorf("Result returned by tag is correct: %v", diff)
	}

	s.mgoSession.DB(DBName).C(CollectionName).RemoveAll(nil)
}
