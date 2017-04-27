package perftest

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateWorker(t *testing.T) {
	m := Create()
	var rp RatingParams
	sc := mockStatsController{}
	tp := TestParams{TestID: "abc"}
	tp.DbController = &sc
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
	m := Create()
	var rp RatingParams
	tp := TestParams{TestID: "abc"}
	tp.AdditionalInfo = map[string]string{
		"p1": "1",
	}
	sc := mockStatsController{}
	tp.DbController = &sc
	rp.TestParams = tp

	w := createWorker(m, &rp)

	go w.run()
	time.Sleep(2 * waitTime)
	w.Request <- struct{}{}
	r := <-w.Response
	if !reflect.DeepEqual(r.GetResult().AdditionalInfo, tp.AdditionalInfo) {
		t.Error("Worker returns correct result")
	}

	w.Exit <- struct{}{}
}
