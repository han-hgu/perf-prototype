package stats

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"

	"github.com/BurntSushi/toml"
	_ "github.com/denisenkom/go-mssqldb"
)

// DBConfig contains the db connection information
// Not useful for now until we accquire multiple DBs
type dbConfig struct {
	Server   string
	Port     int
	UID      string
	Pwd      string
	Database string
}

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

func (c *controller) TearDown() {
	if c.db != nil {
		c.db.Close()
	}
}

// GetUDRRates collects the rates, internally using the eventlog table
// returns the rates collected for now, the number of files processed and the next id you should use
// for next query
func (c *controller) GetUDRRates(filename string, lastEventId uint64) (updatedId uint64, filesProcessed uint, result []float64) {
	updatedId = lastEventId
	filesProcessed = 0
	filesCompletedRxp := regexp.MustCompile("Done Processing File" + ".*" + filename)
	InvalidRatesRxp := regexp.MustCompile("UDRs in 0.0 seconds")
	RateRxp := regexp.MustCompile("([0-9]+)*.([0-9]+)* UDRs/second|([0-9]+)* UDRs/second")
	RateValRxp := regexp.MustCompile("([0-9]+)*.([0-9]+)*")

	q := fmt.Sprintf("select id, result from "+
		"eventlog where id > %v and "+
		"(module = 'UDR Rating' or module = 'UDRRatingEngine') order by id", lastEventId)

	rows, err := c.db.Query(q)
	if err != nil {
		log.Println("WARNING: Stats controller generates an error while getting UDR rates: ", err)
		return updatedId, 0, nil
	}

	var id uint64
	var row string
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &row)
		if err != nil {
			log.Println("WARNING: Stats controller generates an error while scanning a row: ", err)
			return updatedId, 0, nil
		}

		// probably a overkill
		if updatedId < id {
			updatedId = id
		}

		if InvalidRatesRxp.MatchString(row) {
			continue
		}

		if filesCompletedRxp.MatchString(row) {
			filesProcessed++
		}

		if fs := RateRxp.FindString(row); fs != "" {
			fsv := RateValRxp.FindString(fs)

			r, err := strconv.ParseFloat(fsv, 64)
			if err == nil {
				result = append(result, r)
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Println("WARNING: Stats controller generates an error: ", err)
		return updatedId, 0, nil
	}

	return updatedId, filesProcessed, result
}

func (c *controller) GetLastIdFromEventLog(like string) uint64 {
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
		err := rows.Scan(&id)
		fmt.Println("HAN >>>>", err)

		if err != nil {
			return 0
		}
	}

	err = rows.Err()
	if err != nil {
		return 0
	}

	return id
}
