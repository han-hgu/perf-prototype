package perftest

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type mockStatsController struct {
}

func (*mockStatsController) UpdateRatingResult(t *TestInfo, dbIDTracker *DBIDTracker) error {
	return nil
}

func (*mockStatsController) UpdateBillingResult(t *TestInfo, dbIDTracker *DBIDTracker) error {
	return nil
}

func (*mockStatsController) UpdateBaselineIDs(dbIDTracker *DBIDTracker) error {
	return nil
}

func (*mockStatsController) UpdateDBParameters(dbname string, dbp *DBParam) error {
	return nil
}

func TestCreate(t *testing.T) {
	m := Create()

	if len(m.s.info) != 0 {
		t.Error("Create() creates a manager")
	}
}

func TestAdd(t *testing.T) {
	sc := mockStatsController{}
	m := Create()
	tp := RatingParams{}
	tp.DbController = &sc
	m.Add("abc", &tp)

	m.RLock()
	defer m.RUnlock()
	w, ok := m.workerMap["abc"]
	if !ok {
		t.Error("workerMap updated")
	}

	if w.ti.Params.(*RatingParams) != &tp {
		t.Error("worker parameter set successfully")
	}
}

func TestGetInvalidTest(t *testing.T) {
	m := Create()
	if _, e := m.Get("InvalidTestID"); e == nil {
		t.Errorf("Invalid error message received, expect %v, actual %v", errors.New("test doesn't exist"), e)
	}
}

func TestGetValidTestWithWorkerRegistered(t *testing.T) {
	ai := map[string]string{
		"p1": "1",
		"p2": "2",
		"p3": "3",
	}

	tp := TestParams{TestID: "abc"}
	tp.AdditionalInfo = ai

	sc := mockStatsController{}
	rp := RatingParams{}
	rp.TestParams = tp
	rp.DbController = &sc

	m := Create()
	m.Add("abc", &rp)

	r, e := m.Get("abc")
	if e != nil {
		t.Error("Worker returns the test result")
	}

	if !reflect.DeepEqual(r.GetResult().AdditionalInfo, ai) {
		t.Error("Worker returns the correct test result")
	}
}

func TestGetValidTestWithWorkerUnRegistered(t *testing.T) {
	m := Create()

	rre := RatingResult{MinRate: 1,
		AvgRate:        3,
		FilesCompleted: 5}

	tre := TestResult{StartTime: time.Now(),
		Done: false}

	rre.TestResult = tre
	ti := TestInfo{}
	ti.Result = &rre

	m.s.add("abc", &ti)
	tra, e := m.Get("abc")
	if e != nil {
		t.Error("Get test result returns no error")
	}

	if !reflect.DeepEqual(tra, &rre) {
		t.Errorf("Get test result: expect %v, actual %v", &rre, tra)
	}
}
