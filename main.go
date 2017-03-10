package main

import (
	"fmt"
	"regexp"

	"github.com/perf-prototype/stats"
)

func main() {
	defer stats.GetController().TearDown()

	//fmt.Println(stats.GetController().GetLastIdFromEventLog("Running Billing For owner ''Momentum_Voice"))
	fmt.Println(stats.GetController().GetUDRRates("BW", 221354658))

	s := regexp.MustCompile("14.5 UDRs/second")
	fmt.Println(s.FindString("[0-9]+"))

	//re := regexp.MustCompile("([0-9]+)*.([0-9]+)* UDRs/second")
	//fmt.Println(re.FindString("14.5 UDRs/second"))

	//rating.QueryStats("", sc)
}
