package main

import (
	"log"

	"github.com/perf-prototype/stats"
)

func main() {
	log.Println(stats.Conf.Server)
	log.Println(stats.Conf.Port)
	log.Println(stats.Conf.UID)
	log.Println(stats.Conf.Pwd)
	log.Println(stats.Conf.Database)
}
