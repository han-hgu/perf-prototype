package perftest

import (
	"reflect"
	"testing"
)

var s *store

func setup() {
	s = new(store)
}

func TestStoreAddSeperateTests(t *testing.T) {
	setup()
	t1 := new(TestInfo)
	t2 := new(TestInfo)
	t3 := new(TestInfo)
	s.add("t1", t1)
	s.add("t2", t2)
	s.add("t3", t3)

	if len(s.info) != 3 {
		t.Errorf("Corret number of tests added to the store, expect %v, actual %v", len(s.info), 3)
	}
}

func TestStoreAddSameTests(t *testing.T) {
	setup()
	t1 := new(TestInfo)
	s.add("t1", t1)
	e := s.add("t1", t1)

	if e == nil {
		t.Error("Error received if adding an existing testID")
	}
}

func TestStoreGetNonExistTest(t *testing.T) {
	setup()
	if _, e := s.get("non-exisiting"); e == nil {
		t.Error("Error received if if getting an non-existing testID")
	}

}
func TestStoreGetTest(t *testing.T) {
	setup()

	t1 := new(TestInfo)
	p := new(RatingParams)
	p.AmtFieldIndex = 1
	p.TimpstampFieldIndex = 2
	p.NumOfFiles = 3
	p.RawFields = []string{"abc"}
	r := new(RatingResult)
	r.AvgRate = 1.3
	r.FilesCompleted = 2
	t1.Params = p
	t1.Result = r

	s.add("abc", t1)

	tr, e := s.get("abc")
	if e != nil {
		t.Error("No Error if getting an existing testID")
	}

	if !reflect.DeepEqual(*t1, tr) {
		t.Error("Test results are equal")
	}
}

func TestUpdateNonExistTestNilStore(t *testing.T) {
	setup()

	t1 := new(TestInfo)
	if e := s.update("nonexisting", t1); e == nil {
		t.Error("Update a nil store returns error")
	}
}

func TestUpdateNonExistingTest(t *testing.T) {
	setup()

	t1 := new(TestInfo)
	s.add("abc", t1)
	if e := s.update("nonexisting", t1); e == nil {
		t.Error("Update a nil store returns error")
	}
}

func TestUpdateExistingTest(t *testing.T) {
	setup()

	t1 := new(TestInfo)
	t2 := new(TestInfo)
	t2.Params = &RatingParams{AmtFieldIndex: 1}
	s.add("abc", t1)
	if e := s.update("abc", t2); e != nil {
		t.Error("Update an existing test doesn't return error")
	}

	tc, _ := s.get("abc")
	if !reflect.DeepEqual(tc, *t2) {
		t.Error("Test info updated")
	}
}

func TestMain(m *testing.M) {
	m.Run()
}
