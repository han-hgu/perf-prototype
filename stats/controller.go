package stats

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
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
	qUdr := "seelct top 1 id from udr order by id desc"
	return c.getLastID(qUdr)
}

func (c *Controller) getLastUdrExceptionID() uint64 {
	qUdrException := "select top 1 id from udrexception order by id desc"
	return c.getLastID(qUdrException)
}

func (c *Controller) UpdateIDsForRatingTest() {
	//sp.lastEventLogID = c.getLastEventLogID()
	//sp.lastUdrExceptionID = c.getLastUdrExceptionID()
	//sp.lastUdrID = c.getLastUdrID()
}

// GetLastIDFromEventLog to get the last ID from the eventlog table
// @param like
// Used in the like clause in the query
func (c *Controller) GetLastIDFromEventLog(like string) uint64 {
	q := fmt.Sprintf("select top 1 id from "+
		"eventlog where result like '%%%s%%' order by id desc", like)

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
