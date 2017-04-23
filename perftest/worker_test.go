package perftest

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateWorker(t *testing.T) {
	dbc := mockStatsController{}
	m := Create(&dbc)
	var rp RatingParams
	tp := TestParams{TestID: "abc"}
	tp.AdditionalInfo = map[string]string{
		"p1": "1",
	}
	rp.TestParams = tp
	w := createWorker(m, &rp)

	if w.tt != RATING {
		t.Error("Worker type is RATING")
	}
}

func TestRun(t *testing.T) {
	dbc := mockStatsController{}
	m := Create(&dbc)
	var rp RatingParams
	tp := TestParams{TestID: "abc"}
	tp.AdditionalInfo = map[string]string{
		"p1": "1",
	}
	rp.TestParams = tp
	w := createWorker(m, &rp)
	waitTime = 1 * time.Microsecond
	go w.run()
	time.Sleep(2 * waitTime)
	w.Request <- struct{}{}
	r := <-w.Response
	if !reflect.DeepEqual(r.GetResult().AdditionalInfo, tp.AdditionalInfo) {
		t.Error("Worker returns correct result")
	}

	w.Exit <- struct{}{}
}
