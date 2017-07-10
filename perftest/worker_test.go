package perftest

import (
	"reflect"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestCreateWorkerForRating(t *testing.T) {
	testStore := new(mockStore)
	testStore.Initialize()
	m := Create(testStore)
	var rp RatingParams
	sc := mockStatsController{}
	tp := TestParams{ID: bson.NewObjectId()}
	tp.DbController = &sc
	tp.Cmt = "This is a comment."
	rp.TestParams = tp
	w := createWorker(m, &rp)

	if w.tt != RATING {
		t.Error("Worker type is RATING")
	}
}

func TestCreateWorkerForBilling(t *testing.T) {
	testStore := new(mockStore)
	testStore.Initialize()
	m := Create(testStore)
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
	testStore := new(mockStore)
	testStore.Initialize()
	m := Create(testStore)
	var rp RatingParams
	tp := TestParams{ID: bson.NewObjectId()}
	tp.Cmt = "This is a comment."
	sc := mockStatsController{}
	tp.DbController = &sc
	rp.TestParams = tp

	w := createWorker(m, &rp)

	go w.run()
	time.Sleep(2 * waitTime)
	w.Request <- struct{}{}
	r := <-w.Response
	if r.Result().Cmt != tp.Cmt {
		t.Error("Worker returns correct result")
	}

	w.Exit <- struct{}{}
}

func TestWorkerAddResultsToStoreWhenDone(t *testing.T) {
	testStore := new(mockStore)
	testStore.Initialize()
	m := Create(testStore)
	var rp RatingParams
	testID := bson.NewObjectId()
	tp := TestParams{ID: testID}
	tp.Cmt = "This is a comment."
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

	w.tm.workerMap[testID] = w

	go w.tm.Get(testID)
	go w.tm.Get(testID)
	go w.tm.Get(testID)
	trs, e := w.tm.Get(testID)
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
