package stats

import (
	"database/sql"
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

// UpdateDBParameters to update the database parameters
func (c *Controller) UpdateDBParameters(dbname string, dbp *perftest.DBParam) error {
	dbp.CompatibilityLevel = c.compatiblityLevel(dbname)
	return nil
}

func (c *Controller) compatiblityLevel(dbname string) (clevel uint8) {
	q := `SELECT compatibility_level FROM sys.databases WHERE name = '` + dbname + `'`
	c.getLastVal(q, &clevel)
	return clevel
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

func (c *Controller) getLastVal(q string, v interface{}) {
	err := c.db.QueryRow(q).Scan(v)

	if err != nil {
		log.Fatal("ERR:", err)
	}
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
