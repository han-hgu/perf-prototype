package main

import (
	"fmt"
	"os"

	"github.com/perf-prototype/stats"
)

func main() {
	defer stats.GetController().TearDown()

	//fmt.Println(stats.GetController().GetLastIdFromEventLog("Running Billing For owner ''Momentum_Voice"))
	_, _, r := stats.GetController().GetUDRRates("BW", 221354658)
	f, _ := os.Create("perf.chart")
	defer f.Close()

	fmt.Fprintln(f, "ChartType = line")
	fmt.Fprintln(f, "Title = UDR Rate Comparison")
	fmt.Fprintln(f, "SubTitle = Momentum Broadworks Feed")
	fmt.Fprintln(f, "ValueSuffix = records/second")

	f.WriteString("XAxisNumbers = ")
	for index := range r {
		fmt.Fprint(f, index+1)
		fmt.Fprint(f, ", ")
	}
	f.WriteString("\n\n")

	f.WriteString("YAxisText = UDRs/second\n\n")
	f.WriteString("Data|EngageIP = ")
	for _, element := range r {
		fmt.Fprint(f, element)
		fmt.Fprint(f, ", ")
	}
	f.WriteString("\n\n")
	f.Sync()

	//re := regexp.MustCompile("([0-9]+)*.([0-9]+)* UDRs/second")
	//fmt.Println(re.FindString("14.5 UDRs/second"))

	//rating.QueryStats("", sc)
}
