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

// statsHandler
func statsHandler(w http.ResponseWriter, r *http.Request) {
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

// TestRequestHandler sets up the test and returns the test id for
// future query
func TestRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	var rparams perftest.RatingParams
	var bparams perftest.BillingParams
	var testID string
	var e error
	decoder := json.NewDecoder(r.Body)
	switch category {
	case "rating":
		if e = decoder.Decode(&rparams); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		if testID, e = StartRatingTest(&rparams); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}

	case "billing":
		if e = decoder.Decode(&bparams); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		if testID, e = StartBillingTest(&bparams); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Invalid category", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"testID": testID})
}

// AddV1Routes adds version 1 handlers
func AddV1Routes(r *mux.Router) {
	r.HandleFunc("/{category:(?:billing|rating)}/tests/{testID}", statsHandler).Methods("GET")
	r.HandleFunc("/{category:(?:billing|rating)}/tests", TestRequestHandler).Methods("POST")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// TODO: OPTIONS handler
	AddV1Routes(r.PathPrefix("/v1").Subrouter())
	log.Fatal(http.ListenAndServe(":4999", r))
}
