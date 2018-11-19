package stats

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/perf-prototype/perftest"
)

// UpdateRatingResult updates the testInfo and database id tracker
func (c *Controller) UpdateRatingResult(ti *perftest.TestInfo, dbIDTracker *perftest.DBIDTracker) error {
	start := time.Now()
	dbIDTracker.EventLogCurrent = c.getLastEventLogID()
	dbIDTracker.UDRCurrent = c.getLastUdrID()
	dbIDTracker.UDRExceptionCurrent = c.getLastUdrExceptionID()

	rp, ok := ti.Params.(*perftest.RatingParams)
	if !ok {
		log.Fatal("ERR: Failed to cast ti.Params to *RatingParams")
	}

	rr, ok := ti.Result.(*perftest.RatingResult)
	if !ok {
		log.Fatal("ERR: Failed to cast ti.Result to *RatingResult")
	}

	var (
		udrC                   uint64
		udrExceptionC          uint64
		numberOfFilesProcessed uint32
		udrTotal               uint64
	)

	var wg sync.WaitGroup
	// UDRs
	wg.Add(1)
	go c.getUDRCount(&wg, dbIDTracker.UDRLastProcessed, dbIDTracker.UDRCurrent, &udrC)

	// This is only for debug, this number is temporarily used for ending the test since the UDR count from segments seems inaccurate
	if rp.NumOfUDRRecords != 0 {
		wg.Add(1)
		go c.getUDRCount(&wg, dbIDTracker.UDRStarted, dbIDTracker.UDRCurrent, &udrTotal)
	}

	// UDRExceptions
	wg.Add(1)
	go c.getUDRExceptionCount(&wg, dbIDTracker.UDRExceptionLastProcessed, dbIDTracker.UDRExceptionCurrent, &udrExceptionC)

	// Number of rating files processed
	if rp.NumOfFiles != 0 {
		wg.Add(1)
		go c.numOfFileProcessed(&wg, rp.FilenamePrefix, dbIDTracker.EventLogLastProcessed, dbIDTracker.EventLogCurrent, &numberOfFilesProcessed)
	}

	wg.Wait()

	if udrC != udrTotal-rr.UDRProcessed {
		log.Printf("DEBUG: udrC reported is %v, but udrTotal: %v, UDRProcessed: %v", udrC, udrTotal, rr.UDRProcessed)
	}
	UDRPreviousProcessed := rr.UDRProcessed
	rr.UDRProcessed = udrTotal
	// Attention: UDRProcessed is a field for charting purpose, in order to
	// make the interface unified, cast uint64 to float32 assuming the cast will
	// always be successful
	rr.UDRProcessedTrend = append(rr.UDRProcessedTrend, rr.UDRProcessed)
	rr.UDRExceptionProcessed += udrExceptionC
	rr.FilesCompleted += numberOfFilesProcessed
	if (rp.NumOfFiles != 0 && rr.FilesCompleted == rp.NumOfFiles) ||
		(rp.NumOfUDRRecords != 0 && rp.NumOfUDRRecords <= udrTotal) {
		log.Printf("DEBUG: rp.NumOfUDRRecords: %v, udrTotal: %v\n", rp.NumOfUDRRecords, udrTotal)
		rr.Done = true
		duration := start.Sub(rr.StartTime)
		rr.Duration = duration.String()
		rr.AvgRate = float32(float64(rr.UDRProcessed) / duration.Seconds())
	}

	// calculate rates by counting
	currRate := float32(float64(rr.UDRProcessed-UDRPreviousProcessed) / float64(start.Sub(dbIDTracker.TimePrevious).Seconds()))
	if rr.MinRate == 0 || currRate < rr.MinRate {
		rr.MinRate = currRate
	}
	rr.Rates = append(rr.Rates, currRate)

	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventLogCurrent
	dbIDTracker.UDRLastProcessed = dbIDTracker.UDRCurrent
	dbIDTracker.UDRExceptionLastProcessed = dbIDTracker.UDRExceptionCurrent
	// don't set it to now since it should be the time we grab the db table IDs
	// and pick up from there
	dbIDTracker.TimePrevious = start
	return nil
}

// numOfFileProcessed returns the number of UDR files shown completed in the
// eventlog between eventlog ID "firstID" and "lastID"
func (c *Controller) numOfFileProcessed(wg *sync.WaitGroup, filename string, firstID, lastID uint64, filesProcessed *uint32) {
	if wg != nil {
		defer wg.Done()
	}

	filesCompletedRxp := regexp.MustCompile("Done Processing File" + ".*" + filename + ".*")
	q := fmt.Sprintf("select id, result from "+
		"eventlog where id > %v and id <= %v and "+
		"(module = 'UDR Rating' or module = 'UDRRatingEngine') order by id", firstID, lastID)

	rows, err := c.db.Query(q)
	if err != nil {
		log.Fatalf("ERR: Stats controller generates an error getting number of files: %v", err)
	}

	var id uint64
	var row string
	defer rows.Close()
	for rows.Next() {
		rowErr := rows.Scan(&id, &row)
		if rowErr != nil {
			log.Fatalf("ERR: Stats controller generates an error while scanning a row: %v", err)
		}

		if filesCompletedRxp.MatchString(row) {
			(*filesProcessed)++
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("WARNING: Stats controller generates an error: %v", err)
	}
}

func (c *Controller) getUDRCount(wg *sync.WaitGroup, last, current uint64, result *uint64) {
	if wg != nil {
		defer wg.Done()
	}

	q := fmt.Sprintf("select count(*) from udr where id > %v and id <= %v", last, current)
	log.Printf("DEBUG: query for UDR total: %v\n", q)

	c.getLastVal(q, []interface{}{result})
}

func (c *Controller) getUDRExceptionCount(wg *sync.WaitGroup, last, current uint64, result *uint64) {
	if wg != nil {
		defer wg.Done()
	}

	q := fmt.Sprintf("select count(*) from udrException where id > %v and id <= %v", last, current)
	c.getLastVal(q, []interface{}{result})
}

func (c *Controller) getLastUdrID() (id uint64) {
	qUdr := "select top 1 id from udr order by id desc"

	valExists, e := c.getLastVal(qUdr, []interface{}{&id})
	if !valExists || e != nil {
		log.Fatalf("getLastUdrID() gets an error: %v", e)
	}

	return id
}

func (c *Controller) getLastUdrExceptionID() (id uint64) {
	qUdrException := "select top 1 id from udrexception order by id desc"

	valExists, e := c.getLastVal(qUdrException, []interface{}{&id})
	if !valExists || e != nil {
		log.Fatalf("getLastUdrExceptionID() gets an error: %v", e)
	}

	return id
}
