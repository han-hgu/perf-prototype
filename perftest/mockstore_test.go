package perftest

import (
	"reflect"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestStoreAddSeperateTests(t *testing.T) {
	s := new(mockStore)
	s.Initialize()
	defer s.Teardown()
	t1 := new(TestResult)
	t2 := new(TestResult)
	t3 := new(TestResult)
	t1.ID = bson.NewObjectId()
	t2.ID = bson.NewObjectId()
	t3.ID = bson.NewObjectId()
	s.add(t1)
	s.add(t2)
	s.add(t3)

	if len(s.info) != 3 {
		t.Errorf("Corret number of tests added to the store, expect %v, actual %v", len(s.info), 3)
	}
}

func TestStoreAddSameTests(t *testing.T) {
	s := new(mockStore)
	s.Initialize()
	defer s.Teardown()
	t1 := new(TestResult)
	t1.ID = bson.NewObjectId()
	s.add(t1)
	t2 := new(TestResult)
	t2.ID = t1.ID
	s.add(t1)
	e := s.add(t2)

	if e != nil {
		t.Error("Add duplicate test results is a no-op")
	}

	s.RLock()
	defer s.RUnlock()
	if len(s.info) != 1 {
		t.Error("Duplicate test is not added")
	}
}

func TestStoreGetTestFromNilStore(t *testing.T) {
	sn := new(mockStore)
	if _, e := sn.get(bson.NewObjectId()); e == nil {
		t.Error("Error received if get with an non-existent testID")
	}
}

func TestStoreGetTest(t *testing.T) {
	s := new(mockStore)
	s.Initialize()
	defer s.Teardown()

	r := new(RatingResult)
	r.ID = bson.NewObjectId()
	r.AvgRate = 1.3
	r.FilesCompleted = 2

	s.add(r)

	rs, e := s.get(r.ID)
	if e != nil {
		t.Error("No Error if getting an existing testID")
	}

	if !reflect.DeepEqual(rs, r) {
		t.Error("Test results are equal")
	}
}

func TestMain(m *testing.M) {
	waitTime = 1 * time.Microsecond
	m.Run()
}
