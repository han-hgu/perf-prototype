package main

/*
// test mssql db
func main() {
	var statsDBConf stats.DBConfig
	statsDBConf.Server = "192.168.1.47"
	statsDBConf.Port = 1433
	statsDBConf.Database = "EngageIP_NonRevenue_85270"
	statsDBConf.UID = "sa"
	statsDBConf.Pwd = "Q@te$t#1"

	tr := time.Time{}
	fmt.Println("tr", tr)
	sc := stats.CreateController(&statsDBConf)

	q := fmt.Sprintf("select '0001-01-01 00:00:00 +0000 UTC'")
	//q := fmt.Sprintf("CONVERT(datetime, 0)")
	var t time.Time
	fmt.Println("HAN >>> t:", t)
	fmt.Println("HAN >>> bool:", t.IsZero())

	// fmt.Println("HAN >>>> billing started:", sc.BillingStarted("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> billing start time:", sc.BillingStartTime("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> billing finished:", sc.BillingFinished("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> billing finish time:", sc.BillingEndTime("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> invoice rendering started:", sc.InvoiceRenderingStarted("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> invoice rendering start time:", sc.InvoiceRenderingStartTime("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> invoice rendering finished:", sc.InvoiceRenderingFinished("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> invoice rendering finish time:", sc.InvoiceRenderingEndTime("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> bill run finished:", sc.BillrunFinished("Momentum_Retail", 221354332))
	// fmt.Println("HAN >>>> bill run finish time:", sc.BillrunEndTime("Momentum_Retail", 221354332))
}

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan struct{})

	select {
	case c <- struct{}{}:
	case <-time.After(5 * time.Second):
		fmt.Println("HAN >>> times up")
	}

	fmt.Println("OK")
}

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

func main() {
	f, _ := os.Create("C:/tmp/sql2333.txt")
	defer f.Close()

	// open a file
	if file, err := os.Open("C:/tmp/sql.txt"); err == nil {

		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			filesCompletedRxp := regexp.MustCompile("'.*'")
			s := string(filesCompletedRxp.Find(scanner.Bytes()))
			if len(s) > 0 {
				log.Println(s[1 : len(s)-1])
				f.WriteString("SELECT is_disabled FROM sys.indexes WHERE name = '" + s[1:len(s)-1] +
					"'\n")
			}
		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}

}
*/

// This is for generating chart
/*
func main() {
	defer stats.Controller().TearDown()

	//fmt.Println(stats.Controller().GetLastIdFromEventLog("Running Billing For owner ''Momentum_Retail"))
	_, _, r := stats.Controller().GetUDRRates("BW", 221354658)
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

}
*/
