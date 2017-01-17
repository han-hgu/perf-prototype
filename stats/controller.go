package stats

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
	_ "github.com/denisenkom/go-mssqldb"
)

var dbc = `
Server = "192.168.1.47"
Port = 1433
UID = "sa"
Pwd = "Q@te$t#1"
Database = "EngageIP_Revenue"
`

var connStr string

// DBConfig contains the db connection information
type DBConfig struct {
	Server   string
	Port     int
	UID      string
	Pwd      string
	Database string
}

type Controller struct {
	Conf DBConfig
}

// Conf contains the db connection information
var Conf DBConfig

func init() {
	if _, err := toml.Decode(dbc, &Conf); err != nil {
		log.Fatal(err)
	}

	connStr := "server=" + Conf.Server + ";port=" + strconv.Itoa(Conf.Port) + ";" +
		"user id=" + Conf.UID + ";" +
		"password=" + Conf.Pwd + ";" +
		"database=" + Conf.Database + ";"

	if _, err := sql.Open("sqlserver", connStr); err != nil {
		log.Fatal(err)
	}
}

func (c *Controller)
