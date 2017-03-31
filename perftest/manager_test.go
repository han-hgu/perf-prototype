package perftest

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/perf-prototype/stats"
)

func TestGetInfo(t *testing.T) {
	var rp RatingParams
	rp.AdditionalInfo = map[string]string{"foo": "v1", "bar": "v2"}
	if !reflect.DeepEqual(rp.AdditionalInfo, rp.GetInfo()) {
		t.Errorf("GetInfo() returns correct value")
	}
}

func TestGetResult(t *testing.T) {
	rr := RatingResult{MinRate: 0,
		AvgRate:        0,
		FilesCompleted: 0}

	tr := TestResult{StartTime: time.Now(),
		LastLog: 0,
		Done:    false}

	rr.TestResult = &tr

	if rr.GetResult() != &tr {
		t.Error("GetResult() returns the correct value")
	}
}

func TestCreate(t *testing.T) {
	dbc := stats.Controller{}
	m := Create(&dbc)

	if m.db != &dbc {
		t.Error("Create() creates a manager with correct db controller")
	}
}

func TestAdd(t *testing.T) {
	dbc := stats.Controller{}
	m := Create(&dbc)
	ti := TestInfo{}
	m.Add("abc", &ti)
	ti, _ = m.s.get("abc")

	if ti.w == nil {
		t.Error("Worker created")
	}

	ti.w.Exit <- struct{}{}
}

func TestGetInvalidTest(t *testing.T) {
	dbc := stats.Controller{}
	m := Create(&dbc)
	if _, e := m.Get("InvalidTestID"); e == nil {
		t.Errorf("Invalid error message received, expected %v, actual %v", errors.New("test doesn't exist"), e)
	}
}

func TestGetValidTestWithWorkerRegistered(t *testing.T) {
	re := RatingResult{
		MinRate:        2,
		AvgRate:        1,
		FilesCompleted: 3}

	tr := TestResult{StartTime: time.Now(),
		LastLog: 0,
		Done:    true}

	re.TestResult = &tr

	dbc := stats.Controller{}
	m := Create(&dbc)
	ti := TestInfo{}
	ti.Result = &re
	m.Add("abc", &ti)

	ra, _ := m.Get("abc")
	rep := &re

	if !reflect.DeepEqual(rep, ra) {
		t.Errorf("Get test result: expected %v, actual %v", rep, ra)
	}
}

func TestGetValidTestWithWorkerUnRegistered(t *testing.T) {
	dbc := stats.Controller{}
	m := Create(&dbc)

	rre := RatingResult{MinRate: 1,
		AvgRate:        3,
		FilesCompleted: 5}

	tre := TestResult{StartTime: time.Now(),
		LastLog: 0,
		Done:    false}

	rre.TestResult = &tre
	ti := TestInfo{}
	ti.Result = &rre

	m.s.add("abc", &ti)
	tra, e := m.Get("abc")
	if e != nil {
		t.Error("Get test result returns no error")
	}

	if !reflect.DeepEqual(tra, &rre) {
		t.Errorf("Get test result: expected %v, actual %v", &rre, tra)
	}
}
