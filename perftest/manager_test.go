package perftest

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var testStore = new(mockStore)

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

func (*mockStatsController) UpdateDBParameters(dbConf *DBConf, dbp *DBParams) error {
	return nil
}

func (*mockStatsController) TrackKPI(wg *sync.WaitGroup, dbname string, cpu *float32, lr *uint64, lw *uint64, pr *uint64) {
	if wg != nil {
		defer wg.Done()
	}
	return
}

func (*mockStatsController) GetEIPOptions(ac *AppConf) error {
	return nil
}

func TestCreate(t *testing.T) {
	testStore.Initialize()
	m := Create(testStore)

	trs, err := m.GetAll(nil)
	if err != nil || len(trs) != 0 {
		t.Error("Create() creates a manager")
	}
	if m.workerMap == nil {
		t.Error("Manager's worker map initialized")
	}
}

func TestAdd(t *testing.T) {
	sc := mockStatsController{}
	testStore.Initialize()
	m := Create(testStore)
	tp := RatingParams{}
	tp.DbController = &sc
	ntID := bson.NewObjectId()
	m.Add(ntID, &tp)

	m.RLock()
	defer m.RUnlock()
	w, ok := m.workerMap[ntID]
	if !ok {
		t.Error("workerMap updated")
	}

	if w.ti.Params.(*RatingParams) != &tp {
		t.Error("worker parameter set successfully")
	}
}

func TestGetInvalidTest(t *testing.T) {
	testStore.Initialize()
	m := Create(testStore)
	if _, e := m.Get(bson.NewObjectId()); e == nil {
		t.Errorf("Invalid error message received, expect %v, actual %v", errors.New("test doesn't exist"), e)
	}
}

func TestGetValidTestWithWorkerRegistered(t *testing.T) {
	c := "This is a comment"

	tp := TestParams{ID: "abc"}
	tp.Cmt = c

	sc := mockStatsController{}
	rp := RatingParams{}
	rp.TestParams = tp
	rp.DbController = &sc

	testStore.Initialize()
	m := Create(testStore)

	testID := bson.NewObjectId()
	m.Add(testID, &rp)

	r, e := m.Get(testID)
	if e != nil {
		t.Error("Worker returns the test result")
	}

	if r.Result().Cmt != c {
		t.Error("Worker returns the correct test result")
	}
}

func TestGetValidTestWithWorkerUnRegistered(t *testing.T) {
	testStore.Initialize()
	m := Create(testStore)

	rre := RatingResult{MinRate: 1,
		AvgRate:        3,
		FilesCompleted: 5}

	tre := TestResult{}

	tre.StartTime = time.Now()
	tre.Done = false

	rre.TestResult = tre
	testID := bson.NewObjectId()
	rre.ID = testID

	m.s.add(&rre)
	tra, e := m.Get(testID)
	if e != nil {
		t.Error("Get test result returns no error")
	}

	if !reflect.DeepEqual(tra, &rre) {
		t.Errorf("Get test result: expect %v, actual %v", &rre, tra)
	}
}

func TestContains(t *testing.T) {
	var a, b []string

	if !contains(a, b) {
		t.Error("a nil slice is contained in another nil slice")
	}

	b = []string{"abc", "def"}
	if !contains(a, b) {
		t.Error("a nil slice is contained in any non-nil slice")
	}

	a = []string{"x"}
	if contains(a, b) {
		t.Error("slices have different elements so they are not contained")
	}

	a = []string{"abc"}
	if !contains(a, b) {
		t.Error("All elements in slice A is in slice B, contains(A, B) returns true")
	}

	a = []string{"abc", "def"}
	if !contains(a, b) {
		t.Error("contains(A, A) returns true")
	}
}
