package stats

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/perf-prototype/perftest"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
	)

	var wg sync.WaitGroup
	// UDRs
	wg.Add(1)
	go c.getUDRCount(&wg, dbIDTracker.UDRLastProcessed, dbIDTracker.UDRCurrent, &udrC)

	// UDRExceptions
	wg.Add(1)
	go c.getUDRExceptionCount(&wg, dbIDTracker.UDRExceptionLastProcessed, dbIDTracker.UDRExceptionCurrent, &udrExceptionC)

	// Number of rating files processed
	wg.Add(1)
	go c.numOfFileProcessed(&wg, rp.FilenamePrefix, dbIDTracker.EventLogLastProcessed, dbIDTracker.EventLogCurrent, &numberOfFilesProcessed)

	wg.Wait()

	mem, _ := mem.VirtualMemory()
	cpu, _ := cpu.Percent(0, false)
	if mem.UsedPercent > rr.MemMax {
		rr.MemMax = mem.UsedPercent
	}

	if cpu[0] > rr.CPUMax {
		rr.CPUMax = cpu[0]
	}

	fmt.Printf("MemUsedPercent:\t%f%%\n", mem.UsedPercent)
	fmt.Printf("CPUPercent:\t%f%%\n", cpu[0])

	rr.UDRProcessed += udrC
	rr.UDRExceptionProcessed += udrExceptionC
	rr.FilesCompleted += numberOfFilesProcessed
	if rr.FilesCompleted == rp.NumOfFiles {
		rr.Done = true
		duration := start.Sub(rr.StartTime)
		rr.Duration = duration.String()
		rr.AvgRate = float32(float64(rr.UDRProcessed) / duration.Seconds())
	}

	// calculate rates by counting
	currRate := float32(float64(udrC) / float64(start.Sub(dbIDTracker.TimePrevious).Seconds()))
	if rr.MinRate == 0 || currRate < rr.MinRate {
		rr.MinRate = currRate
	}
	rr.Rates = append(rr.Rates, currRate)

	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventLogCurrent
	dbIDTracker.UDRLastProcessed = dbIDTracker.UDRCurrent
	dbIDTracker.UDRExceptionLastProcessed = dbIDTracker.UDRExceptionCurrent
	fmt.Printf("TimeElapsed:\t%v\n\n", time.Since(start))
	// don't set it to now since it should be the time we grab the db table IDs
	// and pick up from there
	dbIDTracker.TimePrevious = start
	return nil
}

func (c *Controller) getRatesFromEventLog(wg *sync.WaitGroup, firstID, lastID uint64, rates *[]float32) {
	if wg != nil {
		defer wg.Done()
	}

	var (
		InvalidRatesRxp = regexp.MustCompile("UDRs in 0.0 seconds")
		RateRxp         = regexp.MustCompile("([0-9]+)*.([0-9]+)* UDRs/second|([0-9]+)* UDRs/second")
		RateValRxp      = regexp.MustCompile("([0-9]+)*.([0-9]+)*")
	)

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

		if InvalidRatesRxp.MatchString(row) {
			continue
		}

		if fs := RateRxp.FindString(row); fs != "" {
			fsv := RateValRxp.FindString(fs)

			r, err2 := strconv.ParseFloat(fsv, 32)
			if err2 == nil {
				*rates = append(*rates, float32(r))
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("WARNING: Stats controller generates an error: %v", err)
	}
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
	*result = c.getRecordCount(q)
}

func (c *Controller) getUDRExceptionCount(wg *sync.WaitGroup, last, current uint64, result *uint64) {
	if wg != nil {
		defer wg.Done()
	}

	q := fmt.Sprintf("select count(*) from udrException where id > %v and id <= %v", last, current)
	*result = c.getRecordCount(q)
}

func (c *Controller) getLastUdrID() uint64 {
	qUdr := "select top 1 id from udr order by id desc"
	return c.getLastID(qUdr)
}

func (c *Controller) getLastUdrExceptionID() uint64 {
	qUdrException := "select top 1 id from udrexception order by id desc"
	return c.getLastID(qUdrException)
}
