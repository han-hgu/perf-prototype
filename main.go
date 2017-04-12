package main

/*
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

}
*/
