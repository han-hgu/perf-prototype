package stats

import (
	"database/sql"
	"log"
	"strconv"
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

// UpdateBaselineIDs updates the IDs of each table to examine
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

// getLastVal finds the first set of values from the query, it returns false
// if it can't get any values
func (c *Controller) getLastVal(q string, v []interface{}) (bool, error) {
	rows, err := c.db.Query(q)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	for rows.Next() {
		rowErr := rows.Scan(v...)
		if rowErr != nil {
			return false, rowErr
		}
	}

	err = rows.Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Controller) getLastEventLogID() (id uint64) {
	qEventLog := "select top 1 id from eventlog order by id desc"
	valExists, e := c.getLastVal(qEventLog, []interface{}{&id})
	if !valExists || e != nil {
		log.Fatalf("getLastEventLogID() gets an error: %v", e)
	}
	return id
}
