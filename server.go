package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perf-prototype/perftest"
)

type httpStats struct {
	Success     uint64
	InvalidBody uint64
}

// ratingStatsHandler
func ratingStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testID := vars["testID"]

	result, e := GetResult(testID)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

// billingStatsHandler
func billingStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("TODO:"))
}

// ratingTestRequestHandler
// ratingTestRequestHandler sets up asrrsssss rating test and returns the test id for
// future query
func ratingTestRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Decode body to rating.controller.testParams object
	var params perftest.RatingParams
	decoder := json.NewDecoder(r.Body)

	if e := decoder.Decode(&params); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	testID, e := StartRatingTest(&params)
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
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// TODO: OPTIONS handler
	AddV1Routes(r.PathPrefix("/v1").Subrouter())
	log.Fatal(http.ListenAndServe(":4999", r))
}
