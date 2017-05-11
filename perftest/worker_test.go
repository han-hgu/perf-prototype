package perftest

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateWorkerForRating(t *testing.T) {
	m := Create()
	var rp RatingParams
	sc := mockStatsController{}
	tp := TestParams{ID: "abc"}
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

func TestCreateWorkerForBilling(t *testing.T) {
	m := Create()
	var bp BillingParams
	sc := mockStatsController{}

	tp := TestParams{}
	tp.DbController = &sc
	bp.TestParams = tp
	w := createWorker(m, &bp)
	if w.tt != BILLING {
		t.Error("Worker type is BILLING")
	}

	go w.run()
	w.Request <- struct{}{}
	<-w.Response
	w.Exit <- struct{}{}
}

func TestRun(t *testing.T) {
	m := Create()
	var rp RatingParams
	tp := TestParams{ID: "abc"}
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
	if !reflect.DeepEqual(r.Result().AdditionalInfo, tp.AdditionalInfo) {
		t.Error("Worker returns correct result")
	}

	w.Exit <- struct{}{}
}

func TestWorkerAddResultsToStoreWhenDone(t *testing.T) {
	m := Create()
	var rp RatingParams
	tp := TestParams{ID: "abc"}
	tp.AdditionalInfo = map[string]string{
		"p1": "1",
	}
	sc := mockStatsController{}
	tp.DbController = &sc
	rp.TestParams = tp

	w := createWorker(m, &rp)
	r := w.ti.Result.(*RatingResult)
	r.Done = true
	w.update()
	go w.run()
	w.Request <- struct{}{}
	<-w.Response

	w.tm.workerMap["abc"] = w

	go w.tm.Get("abc")
	go w.tm.Get("abc")
	go w.tm.Get("abc")
	trs, e := w.tm.Get("abc")
	if e != nil {
		t.Error("Worker saves the result to store when test is done")
	}

	if !reflect.DeepEqual(trs, r) {
		t.Error("Worker returns correct result")
	}

	<-w.Response

	if len(w.tm.workerMap) != 0 {
		t.Error("Worker exits after the test is completed")
	}
}
