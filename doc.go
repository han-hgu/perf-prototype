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
*/

package main
