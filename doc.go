// Copyright information

/*
Logisense performance test server.

The server is responsible for
    * Monitoring and posting perf test results
*/

/* Reflection related
// cast example
result := tinfo.Result.(*perftest.RatingResult)
params := tinfo.Params.(*perftest.RatingParams)

// check example
if _, ok := t.Params.(*RatingParams); ok {
    w := createWorker()
    go w.run(testID, tm)
}

// check type example
if reflect.TypeOf(t) != reflect.TypeOf(s.info[uuid].Result) {
    log.Println("FATAL: testResult type mis-match")
    return errors.New("testResult type mis-match")
}

// EIP version representation
<?xml version="1.0" encoding="utf-8"?>
<Version>
  <Major>8</Major>
  <Minor>5</Minor>
  <Revision>26</Revision>
  <Build>5</Build>
  <DB>66</DB>
  <!-- Hotfix.3 2017-01-13 -->
</Version>

// Sample rating post Request for Broadworks feed
{
	"db_config":{
		"ip": "192.168.1.47",
		"port": 1433,
		"db_name": "EngageIP_NonRevenue_85270",
		"uid": "sa",
		"password": "Q@te$t#1"
	},

	"additional_info":{
      "version":"8.5.26.5-Hotfix.62RC5",
      "cpu_app_server":"Intel(R) Xeon(R) CPU E5-2680 v3 @2.50GHz",
      "mem_app_server":"87.9GB",
      "cpu_db":"Intel(R) Xeon(R) CPU E5-2680 v3 @2.50GHz (2 processors)",
      "mem_db":"256GB",
      "cache_size":"10KB",
      "comment": "Perf result for index change for XXX",
      "notes": "This is for demo"
	},

   "use_existing_file":true,

   "amount_field_index":0,
   "timestamp_field_index":8,
   "number_of_records_per_file":200,
   "number_of_files":50,
   "drop_location":"D:/EngageIP/NonRevenue/UsageData/Broadworks",
   "raw_fields":[
      "someuuih",
      "3100000277",
      "Normal",
      "+17209431929",
      "+17209435003",
      "Terminating",
      "+17209035791",
      "Public",
      "+17209431929",
      "20170202041058.344",
      "1-060000",
      "Yes-PostRedirection",
      "20170202182958.398",
      "20170202185958.309",
      "016",
      "VoIP",
      "",
      "",
      "",
      "",
      "",
      "",
      "remote",
      "4.55.5.33",
      "SD1ffmf01-0adccc7baac1643be8e71c2cc226c928-v3000v3",
      "PCMU/8000",
      "",
      "",
      "",
      "",
      "",
      "3100000276-01",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "y",
      "",
      "",
      "30351523705:0",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "7209431929@mymtm.us",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "No",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "47.910",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "Network",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "No",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "",
      "1172978:1",
      "",
      "",
      "",
      "",
      "",
      ""
   ]
}

// sample billing test
{
	"db_config":{
		"ip": "192.168.1.47",
		"port": 1433,
		"db_name": "EngageIP_NonRevenue_85270",
		"uid": "sa",
		"password": "Q@te$t#1"
	},

	"additional_info":{
      "version":"8.5.26.5-Hotfix.62RC5",
      "comment": "Perf result for index change for XXX"
	},

   "owner_name": "Momentum_Retail"
}

func (c *Controller) getRatesFromEventLog(wg *sync.WaitGroup, firstID, lastID uint64, rates *[]float32) {
	if wg != nil {
		defer wg.Done()
	}

	var (
		InvalidRatesRxp = regexp.MustCompile("UDRs in 0.0 seconds")
		RateRxp         = regexp.MustCompile("([0-9]+)*.([0-9]+)* UDRs/second|([0-9]+)* UDRs/second")
		RateValRxp      = regexp.MustCompile("([0-9]+)*.([0-9]+)*")
	)

	q := fmt.Sprintf("select id, result from "+
		"eventlog where id > %v and id <= %v and "+
		"(module = 'UDR Rating' or module = 'UDRRatingEngine') order by id", firstID, lastID)

	rows, err := c.db.Query(q)
	if err != nil {
		log.Fatalf("ERR: Stats controller generates an error getting number of files: %v", err)
	}

	var id uint64
	var row string
	defer rows.Close()
	for rows.Next() {
		rowErr := rows.Scan(&id, &row)
		if rowErr != nil {
			log.Fatalf("ERR: Stats controller generates an error while scanning a row: %v", err)
		}

		if InvalidRatesRxp.MatchString(row) {
			continue
		}

		if fs := RateRxp.FindString(row); fs != "" {
			fsv := RateValRxp.FindString(fs)

			r, err2 := strconv.ParseFloat(fsv, 32)
			if err2 == nil {
				*rates = append(*rates, float32(r))
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("WARNING: Stats controller generates an error: %v", err)
	}
}

*/

package main
