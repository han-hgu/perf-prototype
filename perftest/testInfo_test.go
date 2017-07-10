package perftest

import (
	"reflect"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestGetComment(t *testing.T) {
	var rp RatingParams
	rp.Cmt = "Rating test commment."
	if rp.Cmt != rp.Comment() {
		t.Errorf("Comment() returns correct value")
	}
}

func TestGetTestID(t *testing.T) {
	var rp RatingParams
	id := bson.NewObjectId()
	tp := TestParams{ID: id}
	rp.TestParams = tp
	if rp.TestID() != id {
		t.Error("TestID() returns correct value")
	}
}

func TestGetTestIDFromResult(t *testing.T) {
	testID := bson.NewObjectId()
	tr := TestResult{}
	tr.ID = testID
	if tr.TestID() != testID {
		t.Error("Get ID from testResult returns correct value")
	}
}

func TestChartTitleFromResult(t *testing.T) {
	testID := bson.NewObjectId()
	tr := TestResult{}

	tr.ID = testID
	// No title in the TestResult
	if tr.ChartTitle() != testID.Hex() {
		t.Error("Test ID is used as the chart title if chart title is empty")
	}

	testTitle := "testTitleTestChartTitleFromResult"
	tr.CTitle = testTitle

	if tr.ChartTitle() != testTitle {
		t.Error("ChartTitle() returns correct value")
	}
}

func TestGetResult(t *testing.T) {
	rr := RatingResult{MinRate: 0,
		AvgRate:        0,
		FilesCompleted: 0}

	tr := TestResult{}

	tr.StartTime = time.Now()
	tr.Done = false

	rr.TestResult = tr

	if !reflect.DeepEqual(rr.Result(), &tr) {
		t.Error("Result() returns correct value")
	}
}

func TestAddAppServerCPU(t *testing.T) {
	var tr TestResult
	tr.AddAppServerCPU(0.1)
	tr.AddAppServerCPU(0.2)
	tr.AddAppServerCPU(0.3)

	if len(tr.AppStats.CPU) != 3 {
		t.Error("CPU samples are added")
	}

	if tr.AppStats.CPUMaxium != 0.3 {
		t.Error("CPUMaxium is updated")
	}
}

func TestAddDBCPU(t *testing.T) {
	var tr TestResult
	tr.AddDBCPU(0.1)
	tr.AddDBCPU(0.2)
	tr.AddDBCPU(0.3)

	if len(tr.DBStats.CPU) != 3 {
		t.Error("CPU samples are added")
	}

	if tr.DBStats.CPUMaxium != 0.3 {
		t.Error("CPUMaxium is updated")
	}
}

func TestAddAppServerMem(t *testing.T) {
	var tr TestResult
	tr.AddAppServerMem(0.1)
	tr.AddAppServerMem(0.2)
	tr.AddAppServerMem(0.3)

	if len(tr.AppStats.Mem) != 3 {
		t.Error("CPU samples are added")
	}

	if tr.AppStats.MemMaxium != 0.3 {
		t.Error("CPUMaxium is updated")
	}
}

func TestAddLogicalReads(t *testing.T) {
	var tr TestResult
	tr.AddLogicalReads(10)
	tr.AddLogicalReads(20)
	tr.AddLogicalReads(30)

	if tr.DBStats.LReadsTotal != 60 {
		t.Error("Correct number of logical reads added")
	}
}

func TestFetchAppServerCPUStats(t *testing.T) {
	var tr RatingResult
	tr.AddAppServerCPU(0.1)
	tr.AddAppServerCPU(0.2)
	tr.AddAppServerCPU(0.3)
	if !reflect.DeepEqual(FetchAppServerCPUStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch app server CPU stats returns the correct result")
	}
}

func TestFetchAppServerMemStats(t *testing.T) {
	var tr RatingResult
	tr.AddAppServerMem(0.1)
	tr.AddAppServerMem(0.2)
	tr.AddAppServerMem(0.3)
	if !reflect.DeepEqual(FetchAppServerMemStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch app server Mem stats returns the correct result")
	}
}

func TestFetchDBServerCPUStats(t *testing.T) {
	var tr BillingResult
	tr.AddDBServerCPU(0.1)
	tr.AddDBServerCPU(0.2)
	tr.AddDBServerCPU(0.3)
	if !reflect.DeepEqual(FetchDBServerCPUStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch DB server CPU stats returns the correct result")
	}
}

func TestFetchDBServerMemStats(t *testing.T) {
	var tr BillingResult
	tr.AddDBServerMem(0.1)
	tr.AddDBServerMem(0.2)
	tr.AddDBServerMem(0.3)
	if !reflect.DeepEqual(FetchDBServerMemStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch app server mem stats returns the correct result")
	}
}

func TestFetchRates(t *testing.T) {
	var rr RatingResult
	rr.Rates = append(rr.Rates, 0.1)
	rr.Rates = append(rr.Rates, 0.2)
	rr.Rates = append(rr.Rates, 0.3)
	if !reflect.DeepEqual(FetchRates(&rr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch rates returns the correct result")
	}
}

func TestFetchUDRProcessedTrend(t *testing.T) {
	var rr RatingResult
	rr.UDRProcessedTrend = append(rr.UDRProcessedTrend, 10000)
	rr.UDRProcessedTrend = append(rr.UDRProcessedTrend, 20000)
	rr.UDRProcessedTrend = append(rr.UDRProcessedTrend, 30000)
	if !reflect.DeepEqual(FetchUDRProcessedTrend(&rr), []uint64{10000, 20000, 30000}) {
		t.Error("Fetch UDR Processed Trend returns the correct result")
	}
}
