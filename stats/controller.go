package stats

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/BurntSushi/toml"
	_ "github.com/denisenkom/go-mssqldb"
)

// dbConfig contains the db connection information
type dbConfig struct {
	Server   string
	Port     int
	UID      string
	Pwd      string
	Database string
}

// unexported type, calling GetController() to get the singleton
type controller struct {
	conf       *dbConfig
	connString string
	db         *sql.DB
}

var c *controller
var once sync.Once

// GetController gets a singleton to communicate with the db, only
// Supports single db for now
func GetController() *controller {
	once.Do(func() {
		var conf dbConfig
		if _, err := toml.DecodeFile("perf.conf", &conf); err != nil {
			log.Fatal(err)
		}

		c = new(controller)
		c.conf = &conf
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
	})

	return c
}

// Teardown to close the database properly
func (c *controller) TearDown() {
	if c.db != nil {
		c.db.Close()
	}
}

// GetLastIDFromEventLog to get the last ID from the eventlog table
// @param like
// Used in the like clause in the query
func (c *controller) GetLastIDFromEventLog(like string) uint64 {
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
		fmt.Println("HAN >>>>", err)

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
