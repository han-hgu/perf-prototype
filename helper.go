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

type templateDataFeedUint64 struct {
	Title   string
	IDs     []string
	Results [][]*uint64
}

type collector func(tr perftest.Result) []float32
type collectorUint64 func(tr perftest.Result) []uint64

func collectSamplesForTemplateUint64(title string, trs []perftest.Result, c collectorUint64) (*templateDataFeedUint64, error) {
	var (
		rv      templateDataFeedUint64
		rawVals [][]uint64
		maxl    int
	)

	rv.Title = title

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.ChartTitle())
		rawVals = append(rawVals, c(tr))
		currLen := len(c(tr))
		if maxl < currLen {
			maxl = currLen
		}
	}

	parseResultsForTemplateDataFeedUint64(&rv, rawVals, maxl)
	return &rv, nil
}

func collectSamplesForTemplate(title string, trs []perftest.Result, c collector) (*templateDataFeed, error) {
	var (
		rv      templateDataFeed
		rawVals [][]float32
		maxl    int
	)

	rv.Title = title

	// index
	rv.IDs = append(rv.IDs, "X")
	for _, tr := range trs {
		rv.IDs = append(rv.IDs, tr.ChartTitle())
		rawVals = append(rawVals, c(tr))
		currLen := len(c(tr))
		if maxl < currLen {
			maxl = currLen
		}
	}

	parseResultsForTemplateDataFeed(&rv, rawVals, maxl)
	return &rv, nil
}

func parseResultsForTemplateDataFeedUint64(tdf *templateDataFeedUint64, rvs [][]uint64, maxlen int) {
	var retVals [][]*uint64
	for i := 0; i < maxlen; i++ {
		retVals = append(retVals, make([]*uint64, 0))
		var index = new(uint64)
		*index = uint64(i)
		retVals[i] = append(retVals[i], index)
		for j := 0; j < len(rvs); j++ {
			val := new(uint64)
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

// AppCPUSamplesForTemplate returns app server CPU stats for google charts
func AppCPUSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	return collectSamplesForTemplate("Application Server CPU utilization", trs, perftest.FetchAppServerCPUStats)
}

// AppMemSamplesForTemplate returns app server memory stats for google charts
func AppMemSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	return collectSamplesForTemplate("Application Server Memory utilization", trs, perftest.FetchAppServerMemStats)
}

// DBCPUSamplesForTemplate returns database CPU stats for google charts
func DBCPUSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	return collectSamplesForTemplate("Database Server CPU utilization", trs, perftest.FetchDBServerCPUStats)
}

// DBLogicalReadsForTemplate returns database logical reads for google charts
func DBLogicalReadsForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("Database Logical Reads", trs, perftest.FetchDBServerLReads)
}

// DBPhysicalReadsForTemplate returns database physical reads for google charts
func DBPhysicalReadsForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("Database Physical Reads", trs, perftest.FetchDBServerPReads)
}

// DBLogicalWrites returns database logical writes for google charts
func DBLogicalWrites(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("Database Logical Writes", trs, perftest.FetchDBServerLWrites)
}

// DBMemSamplesForTemplate returns database memory stats for google charts
func DBMemSamplesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	return collectSamplesForTemplate("Database Server Memory utilization", trs, perftest.FetchDBServerMemStats)
}

// UDRRatesForTemplate returns udr rates for google charts
func UDRRatesForTemplate(trs []perftest.Result) (*templateDataFeed, error) {
	return collectSamplesForTemplate("UDR Rates", trs, perftest.FetchRates)
}

// UDRCurrentProcessedForTemplate returns the total number of UDRs per interval
func UDRCurrentProcessedForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("Total UDR Processed", trs, perftest.FetchUDRProcessedTrend)
}

// InvoiceClosedForTemplate returns the number of invoices closed
func InvoiceClosedForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("Invoices Closed", trs, perftest.FetchInvoicesClosed)
}

// UsageTransactionsGeneratedForTemplate returns the number of usage transactions generated
func UsageTransactionsGeneratedForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("Usage Transactions Generated", trs, perftest.FetchUsageTransactionsGenerated)
}

// MRCTransactionsGeneratedForTemplate returns the number of MRC transactions generated
func MRCTransactionsGeneratedForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("MRC Transactions Generated", trs, perftest.FetchMRCTransactionsGenerated)
}

// BillUDRActionCompletedForTemplate returns the number of BillUDR actions completed
func BillUDRActionCompletedForTemplate(trs []perftest.Result) (*templateDataFeedUint64, error) {
	return collectSamplesForTemplateUint64("BillUDR Actions Completed", trs, perftest.FetchBillUDRActionCompleted)
}

// GetBillingActionChartData prepares the data for comparison bar charts
func GetBillingActionChartData(trs []perftest.Result, BillingActionDurationChartData *[][]interface{}, BillingActionItemCountChartData *[][]interface{}) error {
	*BillingActionDurationChartData = make([][]interface{}, 0)
	*BillingActionItemCountChartData = make([][]interface{}, 0)
	billingActions := make(map[string]struct{}, 0)

	// Create heading row while construct a map for all actions
	firstRowForDuration := make([]interface{}, 0)
	firstRowForItemCount := make([]interface{}, 0)
	firstRowForDuration = append(firstRowForDuration, "Action")
	firstRowForItemCount = append(firstRowForItemCount, "Action")
	for i := 0; i < len(trs); i++ {
		br, _ := trs[i].(*perftest.BillingResult)
		firstRowForDuration = append(firstRowForDuration, br.TestResult.ChartTitle())
		firstRowForItemCount = append(firstRowForItemCount, br.TestResult.ChartTitle())

		for action := range br.ActionDuration {
			if _, ok := billingActions[action]; !ok {
				billingActions[action] = struct{}{}
			}
		}
	}
	(*BillingActionDurationChartData) = append((*BillingActionDurationChartData), firstRowForDuration)
	(*BillingActionItemCountChartData) = append((*BillingActionItemCountChartData), firstRowForItemCount)

	for k := range billingActions {
		rowForDuration := make([]interface{}, 0)
		rowForItemCount := make([]interface{}, 0)
		rowForDuration = append(rowForDuration, k)
		rowForItemCount = append(rowForItemCount, k)
		for i := 0; i < len(trs); i++ {
			br, _ := trs[i].(*perftest.BillingResult)

			if _, ok := br.ActionDuration[k]; ok {
				dur, _ := time.ParseDuration(br.ActionDuration[k]["duration"].(string))
				rowForDuration = append(rowForDuration, dur.Seconds())
				rowForItemCount = append(rowForItemCount, br.ActionDuration[k]["item_count"])
			} else {
				rowForDuration = append(rowForDuration, 0)
				rowForItemCount = append(rowForItemCount, 0)
			}
		}

		(*BillingActionDurationChartData) = append((*BillingActionDurationChartData), rowForDuration)
		(*BillingActionItemCountChartData) = append((*BillingActionItemCountChartData), rowForItemCount)
	}

	return nil
}

// GetBillingActionItemCountChartData prepares the data for comparison bar chart for different billing actions for item count
func GetBillingActionItemCountChartData(trs []perftest.Result) ([][]interface{}, error) {
	return nil, nil
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
