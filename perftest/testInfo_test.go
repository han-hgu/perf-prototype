package perftest

import (
	"reflect"
	"testing"
	"time"
)

func TestGetInfo(t *testing.T) {
	var rp RatingParams
	rp.AdditionalInfo = map[string]string{"foo": "v1", "bar": "v2"}
	if !reflect.DeepEqual(rp.AdditionalInfo, rp.GetInfo()) {
		t.Errorf("GetInfo() returns correct value")
	}
}

func TestGetTestID(t *testing.T) {
	var rp RatingParams
	tp := TestParams{TestID: "abc"}
	rp.TestParams = tp
	if rp.GetTestID() != tp.TestID {
		t.Error("GetTestID() returns correct value")
	}
}

func TestGetResult(t *testing.T) {
	rr := RatingResult{MinRate: 0,
		AvgRate:        0,
		FilesCompleted: 0}

	tr := TestResult{StartTime: time.Now(),
		Done: false}

	rr.TestResult = tr

	if !reflect.DeepEqual(rr.GetResult(), &tr) {
		t.Errorf("GetResult() returns correct value")
	}
}
