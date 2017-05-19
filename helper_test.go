package main

import (
	"fmt"
	"testing"
)

func TestParseResultsForTemplateDataFeed(t *testing.T) {
	dt := templateDataFeed{}
	rvs := [][]float32{{1.1, 1.2, 1.3, 1.4}, {2.1}, {3.1, 3.2}}
	parseResultsForTemplateDataFeed(&dt, rvs, 4)
	fmt.Println("HAN >>>> dt:", dt)
}
