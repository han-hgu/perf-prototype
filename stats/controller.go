package stats

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	// mssql driver
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/perf-prototype/perftest"
)

// DBConfig contains the db connection information
type DBConfig struct {
	Server   string
	Port     int
	UID      string
	Pwd      string
	Database string
}

// Controller to return an instance to communicate with one db instance
type Controller struct {
	conf       *DBConfig
	connString string
	db         *sql.DB
}

var c *Controller
var once sync.Once

// UpdateRatingResult updates the testInfo and database id tracker
func (c *Controller) UpdateRatingResult(ti *perftest.TestInfo, dbIDTracker *perftest.DBIDTracker) error {

	fmt.Println("HAN >>>> before")
	fmt.Println("EventLogLastProcessed:", dbIDTracker.EventLogLastProcessed)
	fmt.Println("EventLogCurrent:", dbIDTracker.EventLogCurrent)
	fmt.Println("EventlogStarted:", dbIDTracker.EventlogStarted)
	fmt.Println("UDRLastProcessed:", dbIDTracker.UDRLastProcessed)
	fmt.Println("UDRCurrent:", dbIDTracker.UDRCurrent)
	fmt.Println("UDRStarted:", dbIDTracker.UDRStarted)
	fmt.Println("UDRExceptionLastProcessed:", dbIDTracker.UDRExceptionLastProcessed)
	fmt.Println("UDRExceptionCurrent:", dbIDTracker.UDRExceptionCurrent)
	fmt.Println("UDRExceptionStarted:", dbIDTracker.UDRExceptionStarted)
	fmt.Println("")

	dbIDTracker.EventLogCurrent = c.getLastEventLogID()
	dbIDTracker.UDRCurrent = c.getLastUdrID()
	dbIDTracker.UDRExceptionCurrent = c.getLastUdrExceptionID()
	timeNow := time.Now()

	rp, ok := ti.Params.(*perftest.RatingParams)
	if !ok {
		log.Fatal("ERR: Failed to cast ti.Params to *RatingParams")
	}

	rr, ok := ti.Result.(*perftest.RatingResult)
	if !ok {
		log.Fatal("ERR: Failed to cast ti.Params to *RatingParams")
	}

	var (
		udrC                   uint64
		udrExceptionC          uint64
		numberOfFilesProcessed uint32
	)

	var wg sync.WaitGroup

	// udrC for UDR
	wg.Add(1)
	go func() {
		defer wg.Done()
		q := fmt.Sprintf("select count(*) from udr where id > %v and id <= %v", dbIDTracker.UDRLastProcessed, dbIDTracker.UDRCurrent)
		udrC = c.getRecordCount(q)
		fmt.Println("HAN >>>> udrC: ", udrC)
	}()

	// udrExceptionC for UDR Exception
	wg.Add(1)
	go func() {
		defer wg.Done()
		q := fmt.Sprintf("select count(*) from udrException where id > %v and id <= %v", dbIDTracker.UDRExceptionLastProcessed, dbIDTracker.UDRExceptionCurrent)
		udrExceptionC = c.getRecordCount(q)
		fmt.Println("HAN >>>> udrExceptionC: ", udrExceptionC)
	}()

	// numberOfFilesProcessed for Number of rating files processed
	wg.Add(1)
	go func() {
		defer wg.Done()
		numberOfFilesProcessed = c.numOfFileProcessed(rp.FilenamePrefix, dbIDTracker.EventLogLastProcessed, dbIDTracker.EventLogCurrent)
		fmt.Println("HAN >>>> numberOfFilesProcessed: ", numberOfFilesProcessed)
	}()

	wg.Wait()

	duration := timeNow.Sub(dbIDTracker.TimePrevious)
	rr.UDRProcessed += udrC
	rr.UDRExceptionProcessed += udrExceptionC
	rate := float32(udrC) / float32(duration.Seconds())
	fmt.Println("HAN >>>> duration:", duration)

	rr.FilesCompleted += numberOfFilesProcessed
	if rr.FilesCompleted == rp.NumOfFiles {
		rr.Done = true
	}

	if rate < rr.MinRate {
		rr.MinRate = rate
	}

	rr.AvgRate = (float32(rr.AvgRate)*float32(len(rr.Rates)) + rate) / float32((len(rr.Rates) + 1))
	rr.Rates = append(rr.Rates, rate)

	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventLogCurrent
	dbIDTracker.UDRLastProcessed = dbIDTracker.UDRCurrent
	dbIDTracker.UDRExceptionLastProcessed = dbIDTracker.UDRExceptionCurrent
	dbIDTracker.TimePrevious = timeNow

	fmt.Println("HAN >>> after")
	fmt.Println("EventLogLastProcessed:", dbIDTracker.EventLogLastProcessed)
	fmt.Println("EventLogCurrent:", dbIDTracker.EventLogCurrent)
	fmt.Println("EventlogStarted:", dbIDTracker.EventlogStarted)
	fmt.Println("UDRLastProcessed:", dbIDTracker.UDRLastProcessed)
	fmt.Println("UDRCurrent:", dbIDTracker.UDRCurrent)
	fmt.Println("UDRStarted:", dbIDTracker.UDRStarted)
	fmt.Println("UDRExceptionLastProcessed:", dbIDTracker.UDRExceptionLastProcessed)
	fmt.Println("UDRExceptionCurrent:", dbIDTracker.UDRExceptionCurrent)
	fmt.Println("UDRExceptionStarted:", dbIDTracker.UDRExceptionStarted)

	return nil
}

// UpdateBaselineIDs updates the last IDs for the related tables so that we start
// examine the rows after those IDs
func (c *Controller) UpdateBaselineIDs(dbIDTracker *perftest.DBIDTracker) error {
	dbIDTracker.EventlogStarted = c.getLastEventLogID()
	dbIDTracker.UDRStarted = c.getLastUdrID()
	dbIDTracker.UDRExceptionStarted = c.getLastUdrExceptionID()
	dbIDTracker.TimePrevious = time.Now()
	// Don't call getLast...() again, logs are advancing at the same time
	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventlogStarted
	dbIDTracker.UDRLastProcessed = dbIDTracker.UDRStarted
	dbIDTracker.UDRExceptionLastProcessed = dbIDTracker.UDRExceptionStarted
	dbIDTracker.EventLogCurrent = dbIDTracker.EventlogStarted
	dbIDTracker.UDRCurrent = dbIDTracker.UDRStarted
	dbIDTracker.UDRExceptionCurrent = dbIDTracker.UDRExceptionStarted

	fmt.Println("HAN >>> baseline:", dbIDTracker)
	return nil
}

// CreateController returns a controller to communicate with the sql db based
// on DBConfig
func CreateController(dbc *DBConfig) *Controller {
	c := new(Controller)
	c.conf = dbc
	c.connString = "server=" + c.conf.Server +
		";port=" + strconv.Itoa(c.conf.Port) + ";" +
		"user id=" + c.conf.UID + ";" +
		"password=" + c.conf.Pwd + ";" +
		"database=" + c.conf.Database + ";"

	db, err := sql.Open("sqlserver", c.connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	c.db = db
	return c
}

// TearDown to close the database properly
func (c *Controller) TearDown() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *Controller) getRecordCount(q string) (count uint64) {
	if err := c.db.QueryRow(q).Scan(&count); err != nil {
		log.Fatalf("ERR: Fail to execute %v, error: %v", q, err)
	}
	return count
}

func (c *Controller) getLastID(q string) uint64 {
	rows, err := c.db.Query(q)
	if err != nil {
		// start from 0
		return 0
	}

	var id uint64
	defer rows.Close()
	for rows.Next() {
		rowErr := rows.Scan(&id)
		if rowErr != nil {
			return 0
		}
	}

	err = rows.Err()
	if err != nil {
		return 0
	}

	return id
}

func (c *Controller) getLastEventLogID() uint64 {
	qEventLog := "select top 1 id from eventlog order by id desc"
	return c.getLastID(qEventLog)
}

func (c *Controller) getLastUdrID() uint64 {
	qUdr := "select top 1 id from udr order by id desc"
	return c.getLastID(qUdr)
}

func (c *Controller) getLastUdrExceptionID() uint64 {
	qUdrException := "select top 1 id from udrexception order by id desc"
	return c.getLastID(qUdrException)
}
