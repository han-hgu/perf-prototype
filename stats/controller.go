package stats

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

var dbc = `
Server = "192.168.1.47"
Port = 1433
UID = "sa"
Pwd = "Q@te$t#1"
Database = "EngageIP_Revenue"
`

// DBConfig contains the db connection information
type DBConfig struct {
	Server   string
	Port     int
	UID      string
	Pwd      string
	Database string
}

type controller struct {
	conf       *DBConfig
	connString string
	db         *sql.DB
}

// New to new a stats controller
func New(dbc *DBConfig) *controller {
	c := new(controller)
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

	c.db = db

	return c
}

func (c *controller) TearDown() {
	if c.db != nil {
		c.db.Close()
	}
}

// func init() {
// 	if _, err := toml.Decode(dbc, &Conf); err != nil {
// 		log.Fatal(err)
// 	}
//
// }
