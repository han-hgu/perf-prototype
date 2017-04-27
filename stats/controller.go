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
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

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
	start := time.Now()
	dbIDTracker.EventLogCurrent = c.getLastEventLogID()
	dbIDTracker.UDRCurrent = c.getLastUdrID()
	dbIDTracker.UDRExceptionCurrent = c.getLastUdrExceptionID()

	fmt.Println("HAN >>>>")
	fmt.Println("EventLogLastProcessed:", dbIDTracker.EventLogLastProcessed)
	fmt.Println("EventLogCurrent:", dbIDTracker.EventLogCurrent)
	fmt.Println("UDRLastProcessed:", dbIDTracker.UDRLastProcessed)
	fmt.Println("UDRCurrent:", dbIDTracker.UDRCurrent)
	fmt.Println("UDRExceptionLastProcessed:", dbIDTracker.UDRExceptionLastProcessed)
	fmt.Println("UDRExceptionCurrent:", dbIDTracker.UDRExceptionCurrent)
	fmt.Println("")

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
		rates                  []float32
		numberOfFilesProcessed uint32
	)

	var wg sync.WaitGroup
	// UDRs
	wg.Add(1)
	go c.getUDRCount(&wg, dbIDTracker.UDRLastProcessed, dbIDTracker.UDRCurrent, &udrC)

	// UDRExceptions
	wg.Add(1)
	go c.getUDRExceptionCount(&wg, dbIDTracker.UDRLastProcessed, dbIDTracker.UDRCurrent, &udrExceptionC)

	// rates
	wg.Add(1)
	go func() {
		defer wg.Done()
		rates = c.getRatesFromEventLog(dbIDTracker.EventLogLastProcessed, dbIDTracker.EventLogCurrent)
	}()

	// Number of rating files processed
	wg.Add(1)
	go func() {
		defer wg.Done()
		numberOfFilesProcessed = c.numOfFileProcessed(rp.FilenamePrefix, dbIDTracker.EventLogLastProcessed, dbIDTracker.EventLogCurrent)
	}()

	wg.Wait()

	// HAN >>>>
	mem, _ := mem.VirtualMemory()
	cpu, _ := cpu.Percent(0, false)
	fmt.Printf("MemUsedPercent:%f%%\n", mem.UsedPercent)
	fmt.Printf("CPUPercent:%f%%\n", cpu[0])

	rr.UDRProcessed += udrC
	rr.UDRExceptionProcessed += udrExceptionC
	rr.FilesCompleted += numberOfFilesProcessed
	fmt.Println("HAN >>> rr.UDRProcessed", rr.UDRProcessed)
	fmt.Println("HAN >>> rr.UDRExceptionProcessed", rr.UDRExceptionProcessed)
	fmt.Println("HAN >>> rr.FilesCompleted", rr.FilesCompleted)
	if rr.FilesCompleted == rp.NumOfFiles {
		rr.Done = true
	}

	for _, v := range rates {
		if v < rr.MinRate && v != 0 {
			rr.MinRate = v
		}

		rr.Rates = append(rr.Rates, v)
	}

	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventLogCurrent
	dbIDTracker.UDRLastProcessed = dbIDTracker.UDRCurrent
	dbIDTracker.UDRExceptionLastProcessed = dbIDTracker.UDRExceptionCurrent
	dbIDTracker.TimePrevious = start
	elapsed := time.Since(start)
	fmt.Println("HAN >>>> Time elapsed:", elapsed)
	fmt.Println("")
	fmt.Println("")
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
