package perftest

import (
	"reflect"
	"testing"
	"time"
)

func TestGetInfo(t *testing.T) {
	var rp RatingParams
	rp.AdditionalInfo = map[string]string{"foo": "v1", "bar": "v2"}
	if !reflect.DeepEqual(rp.AdditionalInfo, rp.Info()) {
		t.Errorf("Info() returns correct value")
	}
}

func TestGetTestID(t *testing.T) {
	var rp RatingParams
	tp := TestParams{ID: "abc"}
	rp.TestParams = tp
	if rp.TestID() != tp.ID {
		t.Error("TestID() returns correct value")
	}
}

func TestGetTestIDFromResult(t *testing.T) {
	testID := "TestGetTestIDFromResult"
	tr := TestResult{ID: testID}
	if tr.TestID() != testID {
		t.Error("Get ID from testResult returns correct value")
	}
}

func TestChartTitleFromResult(t *testing.T) {
	testID := "testIDTestChartTitleFromResult"
	tr := TestResult{
		ID: testID,
	}

	// No title in the TestResult
	if tr.ChartTitle() != testID {
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

	tr := TestResult{StartTime: time.Now(),
		Done: false}

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

func TestFetchDBLReads(t *testing.T) {
	var tr TestResult
	tr.AddLogicalReads(10)
	tr.AddLogicalReads(20)
	tr.AddLogicalReads(30)
	if !reflect.DeepEqual(tr.FetchDBLReads(), []uint64{10, 20, 30}) {
		t.Error("Logical read samples are added")
	}

	if !reflect.DeepEqual(FetchDBServerLReads(&tr), []uint64{10, 20, 30}) {
		t.Error("Logical reads fetched")
	}
}

func TestFetchDBLWrites(t *testing.T) {
	var tr TestResult
	tr.AddLogicalWrites(10)
	tr.AddLogicalWrites(20)
	tr.AddLogicalWrites(30)
	if !reflect.DeepEqual(tr.FetchDBLWrites(), []uint64{10, 20, 30}) {
		t.Error("Logical read samples are added")
	}

	if !reflect.DeepEqual(FetchDBServerLWrites(&tr), []uint64{10, 20, 30}) {
		t.Error("Logical writes fetched")
	}
}

func TestFetchDBPReads(t *testing.T) {
	var tr TestResult
	tr.AddPhysicalReads(10)
	tr.AddPhysicalReads(20)
	tr.AddPhysicalReads(30)
	if !reflect.DeepEqual(tr.FetchDBPReads(), []uint64{10, 20, 30}) {
		t.Error("Logical read samples are added")
	}

	if !reflect.DeepEqual(FetchDBServerPReads(&tr), []uint64{10, 20, 30}) {
		t.Error("Physical read fetched")
	}
}

func TestFetchAppServerCPUStats(t *testing.T) {
	var tr TestResult
	tr.AddAppServerCPU(0.1)
	tr.AddAppServerCPU(0.2)
	tr.AddAppServerCPU(0.3)
	if !reflect.DeepEqual(FetchAppServerCPUStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch app server CPU stats returns the correct result")
	}
}

func TestFetchAppServerMemStats(t *testing.T) {
	var tr TestResult
	tr.AddAppServerMem(0.1)
	tr.AddAppServerMem(0.2)
	tr.AddAppServerMem(0.3)
	if !reflect.DeepEqual(FetchAppServerMemStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch app server Mem stats returns the correct result")
	}
}

func TestFetchDBServerCPUStats(t *testing.T) {
	var tr TestResult
	tr.AddDBServerCPU(0.1)
	tr.AddDBServerCPU(0.2)
	tr.AddDBServerCPU(0.3)
	if !reflect.DeepEqual(FetchDBServerCPUStats(&tr), []float32{0.1, 0.2, 0.3}) {
		t.Error("Fetch DB server CPU stats returns the correct result")
	}
}

func TestFetchDBServerMemStats(t *testing.T) {
	var tr TestResult
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
	if !reflect.DeepEqual(FetchUDRProcessedTrend(&rr), []float32{10000, 20000, 30000}) {
		t.Error("Fetch UDR Processed Trend returns the correct result")
	}
}
