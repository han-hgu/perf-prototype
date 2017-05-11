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

func TestGetResult(t *testing.T) {
	rr := RatingResult{MinRate: 0,
		AvgRate:        0,
		FilesCompleted: 0}

	tr := TestResult{StartTime: time.Now(),
		Done: false}

	rr.TestResult = tr

	if !reflect.DeepEqual(rr.Result(), &tr) {
		t.Errorf("Result() returns correct value")
	}
}
