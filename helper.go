package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/perf-prototype/perftest"
)

// TemplateDataFeed for the data for graphing the udr rate graph
type templateDataFeed struct {
	Title   string
	IDs     []string
	Results [][]*float32
}

func getAppCPUSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	var (
		rv      templateDataFeed
		rawVals [][]float32
		maxl    int
	)

	// title
	rv.Title = "Application Server CPU utilization"

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.TestID())
		rawVals = append(rawVals, tr.AppServerStats().CPU)
		currLen := len(tr.AppServerStats().CPU)
		if maxl < currLen {
			maxl = currLen
		}
	}

	parseResultsForTemplateDataFeed(&rv, rawVals, maxl)

	return &rv, nil
}

func getAppMemSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	var (
		rv      templateDataFeed
		rawVals [][]float32
		maxl    int
	)

	// title
	rv.Title = "Application Server Memory utilization"

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.TestID())
		rawVals = append(rawVals, tr.AppServerStats().Mem)
		currLen := len(tr.AppServerStats().Mem)
		if maxl < currLen {
			maxl = currLen
		}
	}

	parseResultsForTemplateDataFeed(&rv, rawVals, maxl)

	return &rv, nil
}

func getDBCPUSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	var (
		rv      templateDataFeed
		rawVals [][]float32
		maxl    int
	)

	// title
	rv.Title = "Database Server CPU utilization"

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.TestID())
		rawVals = append(rawVals, tr.DBServerStats().CPU)
		currLen := len(tr.DBServerStats().CPU)
		if maxl < currLen {
			maxl = currLen
		}
	}

	parseResultsForTemplateDataFeed(&rv, rawVals, maxl)

	return &rv, nil
}

func getDBMemSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	var (
		rv      templateDataFeed
		rawVals [][]float32
		maxl    int
	)

	// title
	rv.Title = "Database Server Memory utilization"

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.TestID())
		rawVals = append(rawVals, tr.DBServerStats().Mem)
		currLen := len(tr.DBServerStats().Mem)
		if maxl < currLen {
			maxl = currLen
		}
	}

	parseResultsForTemplateDataFeed(&rv, rawVals, maxl)

	return &rv, nil
}

func getUDRRatesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	var (
		rv      templateDataFeed
		rawVals [][]float32
		maxl    int
	)

	// title
	rv.Title = "UDR Rates"

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.TestID())

		rr, ok := tr.(*perftest.RatingResult)
		if !ok {
			return nil, fmt.Errorf("test with id: %v is not a rating test", tr.TestID())
		}

		rawVals = append(rawVals, rr.Rates)
		if maxl < len(rr.Rates) {
			maxl = len(rr.Rates)
		}
	}

	parseResultsForTemplateDataFeed(&rv, rawVals, maxl)

	return &rv, nil
}

func parseResultsForTemplateDataFeed(tdf *templateDataFeed, rvs [][]float32, maxlen int) {
	var retVals [][]*float32
	for i := 0; i < maxlen; i++ {
		retVals = append(retVals, make([]*float32, 0))
		var index = new(float32)
		*index = float32(i)
		retVals[i] = append(retVals[i], index)
		for j := 0; j < len(rvs); j++ {
			val := new(float32)
			if i < len(rvs[j]) {
				*val = rvs[j][i]
			} else {
				val = nil
			}

			retVals[i] = append(retVals[i], val)
		}
	}

	tdf.Results = retVals
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// exists returns true if file path exists
func exists(path string) error {
	// TODO: there are other errors that could be returned not just file doesn't
	// exist one
	_, err := os.Stat(path)
	return err
}

// createFile to create the UDR input files based on the testParams obj
func createFile(t *perftest.RatingParams) error {
	if len(t.RawFields) == 0 {
		return errors.New("Raw fields cannot be empty")
	}

	// check to see if the location exist, location specified must exist
	if err := exists(t.DropLocation); err != nil {
		return err
	}

	var filename string
	for i := uint32(0); i < t.NumOfFiles; i++ {
		filename = t.DropLocation + "/" + t.FilenamePrefix + "-" + strconv.FormatUint(uint64(i), 10) + ".csv"

		fo, err := os.Create(filename)
		defer func() {
			if e := fo.Close(); e != nil {
				panic(e)
			}
		}()

		if err != nil {
			return err
		}

		for i := 0; i < t.NumRecordsPerFile; i++ {
			// No random, rate repeatly using the current timestamp for phase 1
			// 20060102150405 is const have to specify it this way, refer to
			// http://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
			tns := time.Now().Format("20060102150405.000")

			// replace the timestamp
			t.RawFields[t.TimpstampFieldIndex] = tns
			// replace the uniqueness identifier
			var err error
			t.RawFields[0], err = newUUID()
			if err != nil {
				panic(err)
			}
			fo.WriteString(strings.Join(t.RawFields, ",") + "\n")
		}
	}

	return nil
}
