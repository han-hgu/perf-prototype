package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/perf-prototype/perftest"
)

type httpStats struct {
	Success     uint64
	InvalidBody uint64
}

func metaDataRetriever(w http.ResponseWriter, r *http.Request) {
	//https://golangcode.com/get-a-url-parameter-from-a-request/
	tags, _ := r.URL.Query()["tag"]

	mds, e := TestResultSVs(tags)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mds)
}

// statsHandler
func statsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testID := vars["testID"]

	result, e := Result(testID)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func chartHandler(w http.ResponseWriter, r *http.Request) {
	qps := r.URL.Query()
	ids, ok := qps["id"]
	if !ok {
		http.Error(w, "Missing test IDs to compare", http.StatusBadRequest)
		return
	}

	var df struct {
		TestType           string
		CollectionInterval string
		UDRRate            *templateDataFeed
		AppCPU             *templateDataFeed
		AppMem             *templateDataFeed
		DBCPU              *templateDataFeed
		UDRAbsolute        *templateDataFeedUint64
		LReads             *templateDataFeedUint64
		LWrites            *templateDataFeedUint64
		PReads             *templateDataFeedUint64
	}

	var trs []perftest.Result
	var ci string
	for _, v := range ids {
		tr, e := Result(v)

		if e != nil {
			http.Error(w, fmt.Sprintf("Invalid test ID: %s", v), http.StatusBadRequest)
			return
		}

		switch tr.(type) {
		case *perftest.RatingResult:
			if df.TestType == "" {
				df.TestType = "rating"
			} else if df.TestType != "rating" {
				http.Error(w, fmt.Sprintf("Cannot compare tests with different types"), http.StatusBadRequest)
				return
			}

		case *perftest.BillingResult:
			if df.TestType == "" {
				df.TestType = "billing"
			} else if df.TestType != "billing" {
				http.Error(w, fmt.Sprintf("Cannot compare tests with different types"), http.StatusBadRequest)
				return
			}

		default:
			http.Error(w, fmt.Sprintf("Unknown test type: %v", reflect.TypeOf(tr)), http.StatusBadRequest)
			return
		}

		if ci == "" {
			ci = tr.CollectionInterval()
		} else if tr.CollectionInterval() != ci {
			http.Error(w, "Unable to draw comparison graph, test runs must have the same collection intervals", http.StatusBadRequest)
		}

		trs = append(trs, tr)
	}

	df.CollectionInterval = ci
	df.AppCPU, _ = AppCPUSamplesForTemplate(trs)
	df.AppMem, _ = AppMemSamplesForTemplate(trs)
	df.DBCPU, _ = DBCPUSamplesForTemplate(trs)
	df.LReads, _ = DBLogicalReadsForTemplate(trs)
	df.LWrites, _ = DBLogicalWrites(trs)
	df.PReads, _ = DBPhysicalReadsForTemplate(trs)

	if df.TestType == "rating" {
		df.UDRRate, _ = UDRRatesForTemplate(trs)
		df.UDRAbsolute, _ = UDRCurrentProcessedForTemplate(trs)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := template.Must(template.New("comparison.tmpl").ParseFiles("templates/comparison.tmpl")).Execute(w, df); err != nil {
		log.Printf("ERR: html template returns error: %v\n", err)
	}
}

// testRequestHandler sets up the test and returns the test id for
// future query
func testRequestHandler(w http.ResponseWriter, r *http.Request) {
	testType := r.URL.Query().Get("type")

	var rparams perftest.RatingParams
	var bparams perftest.BillingParams
	var testID string
	var e error
	decoder := json.NewDecoder(r.Body)
	switch testType {
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
		http.Error(w, fmt.Sprintf("Invalid/Missing test type, valid types are 'billing' and 'rating', got %v", testType), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": testID})
}

// AddV1Routes adds version 1 handlers
func AddV1Routes(r *mux.Router) {
	r.HandleFunc("/tests", metaDataRetriever).Methods("GET")
	r.HandleFunc("/tests", testRequestHandler).Methods("POST")
	r.HandleFunc("/tests/{testID}", statsHandler).Methods("GET")
	r.HandleFunc("/charts", chartHandler).Methods("GET")
}

func main() {
	defer Teardown()
	r := mux.NewRouter().StrictSlash(true)

	// TODO: OPTIONS handler
	AddV1Routes(r.PathPrefix("/v1").Subrouter())
	log.Fatal(http.ListenAndServe(":4999", r))
}
