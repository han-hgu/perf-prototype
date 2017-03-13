package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perf-prototype/rating"
)

type httpStats struct {
	Success     uint64
	InvalidBody uint64
}

// ratingStatsHandler
func ratingStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testID := vars["testID"]
	fmt.Println(testID)

	// get the db configuration, hard-code to perf.conf
	/*
		var conf stats.DBConfig
		if _, err := toml.DecodeFile("perf.conf", &conf); err != nil {
			log.Fatal(err)
		}

		sc := stats.New(&conf)
		defer sc.TearDown()

		// HAN >>>
		rating.QueryStats(testId, sc)
	*/

}

// billingStatsHandler
func billingStatsHandler(w http.ResponseWriter, r *http.Request) {
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

// ratingTestRequestHandler
// ratingTestRequestHandler sets up a rating test and returns the test id for
// future query
// request format example:
//{
// 	"amount_field_index": 4,
// 	"timestamp_field_index": 1,
// 	"batch_size": 1000,
// 	"number_of_files": 2,
// 	"drop_location": "C:/UsageData",
// 	"raw_fields": [
// 		"7804347632",
// 		"01/14/2009 17:32",
// 		"15196581111",
// 		"7804347632",
// 		"70",
// 		"I",
// 		"value"
// 	]
// }
func ratingTestRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Decode json body to rating.controller.testParams obj
	var params rating.TestParams
	//a, _ := ioutil.ReadAll(r.Body)
	//log.Println(string(a))
	decoder := json.NewDecoder(r.Body)

	if e := decoder.Decode(&params); e != nil {
		//TODO: enum in central location
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	testID, e := rating.StartTest(&params)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"testID": testID})
}

// AddV1Routes adds version 1 handlers
func AddV1Routes(r *mux.Router) {
	r.HandleFunc("/rating/tests/{testID}", ratingStatsHandler).Methods("GET")
	r.HandleFunc("/rating/tests", ratingTestRequestHandler).Methods("POST")
	r.HandleFunc("/billing/tests", billingStatsHandler).Methods("GET")
	r.HandleFunc("/secret", secretHandler).Methods("GET")
}

/*
func main() {
	r := mux.NewRouter().StrictSlash(true)

	// TODO: OPTIONS handler
	AddV1Routes(router.PathPrefix("/v1").Subrouter())
	log.Fatal(http.ListenAndServe(":4999", r))
}
*/
