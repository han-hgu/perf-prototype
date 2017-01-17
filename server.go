package main

import (
	"encoding/json"
	"net/http"

	"github.com/perf-prototype/rating"
)

type httpStats struct {
	Success     uint64
	InvalidBody uint64
}

// ratingStatsHandler
func ratingStatsHandler(w http.ResponseWriter, r *http.Request) {

}

// billingStatsHandler
func billingStatsHandler(w http.ResponseWriter, r *http.Request) {

}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

// udrHandler
// udrHandler responded with a unique identifier after below
// request is received, it creates a record for rating speed and
// update it on demand
// request format:
//{
// 	"amount_field_index": 5,
// 	"timestamp_field_index": 2,
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
func udrHandler(w http.ResponseWriter, r *http.Request) {
	// Decode json body to the rating controller testParams object.
	var params rating.TestParams
	//a, _ := ioutil.ReadAll(r.Body)
	//log.Println(string(a))
	decoder := json.NewDecoder(r.Body)

	if e := decoder.Decode(&params); e != nil {
		//TODO: enum in central location
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	if _, e := rating.StartProcess(&params); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

// func main() {
// 	r := mux.NewRouter()
// 	// TODO: OPTIONS handler etc.
// 	r.HandleFunc("/rating", ratingStatsHandler).Methods("GET")
// 	r.HandleFunc("/rating", udrHandler).Methods("POST")
// 	r.HandleFunc("/billing", billingStatsHandler).Methods("GET")
// 	r.HandleFunc("/secret", secretHandler).Methods("GET")
//
// 	log.Fatal(http.ListenAndServe(":4999", r))
// }
