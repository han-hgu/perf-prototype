package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

	result, e := Result(testID)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func ratingComparisonHandler(w http.ResponseWriter, r *http.Request) {
	qps := r.URL.Query()
	ids, ok := qps["id"]
	if !ok {
		http.Error(w, "Missing test IDs to compare", http.StatusBadRequest)
	}

	var trs []perftest.Result
	for _, v := range ids {
		tr, e := Result(v)
		if e != nil {
			http.Error(w, fmt.Sprintf("Invalid test ID: %s", v), http.StatusBadRequest)
		}

		trs = append(trs, tr)
	}

	ratingDataFeed, e1 := getUDRRatesForTemplate(trs)
	appCPUDataFeed, e2 := getAppCPUSamplesForTemplate(trs)

	if e1 != nil {
		http.Error(w, e1.Error(), http.StatusBadRequest)
	}
	if e2 != nil {
		http.Error(w, e2.Error(), http.StatusBadRequest)
	}

	var df struct {
		UDR    *templateDataFeed
		AppCPU *templateDataFeed
	}

	df.UDR = ratingDataFeed
	df.AppCPU = appCPUDataFeed

	fmt.Println("HAN >>>>> df:", df)

	w.Header().Set("Content-Type", "text/html")
	err := template.Must(template.New("comparison.tmpl").ParseFiles("templates/comparison.tmpl")).Execute(w, df)
	fmt.Println("HAN >>>> e: ", err)

	/*
		var trs []perftest.Result
		for _, v := range ids {
			tr, e := Result(v)
			if e != nil {
				http.Error(w, fmt.Sprintf("Invalid test ID: %s", v), http.StatusBadRequest)
			}

			trs = append(trs, tr)
		}

		// right now only implements rate comparison

		rmax := struct {
			IDs     []string
			Results [][]*float32
		}{
			IDs: []string{"X"},
		}

		rmax.IDs = append(rmax.IDs, ids...)
		for i := 0; i < 3; i++ {
			val := []*float32{}
			var v1 = float32(i)
			var v2 float32 = 1.1

			val = append(val, &v1, &v2, nil)
			rmax.Results = append(rmax.Results, val)

		}
		fmt.Println("HAN >>> rmax.Results", rmax.Results)

		w.Header().Set("Content-Type", "text/html")
		err := template.Must(template.New("comparison.tmpl").ParseFiles("templates/comparison.tmpl")).Execute(w, rmax)
		fmt.Println("HAN >>>> e2", err)
	*/
}

// testRequestHandler sets up the test and returns the test id for
// future query
func testRequestHandler(w http.ResponseWriter, r *http.Request) {
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
	r.HandleFunc("/{category:(?:billing|rating)}/tests", testRequestHandler).Methods("POST")
	r.HandleFunc("/{category:(?:billing|rating)}/tests", ratingComparisonHandler).Methods("GET")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// TODO: OPTIONS handler
	AddV1Routes(r.PathPrefix("/v1").Subrouter())
	log.Fatal(http.ListenAndServe(":4999", r))
}
