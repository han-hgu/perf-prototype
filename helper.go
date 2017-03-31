package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/perf-prototype/perftest"
)

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// exists returns true if file path exists
func exists(path string) error {
	// TODO: there are other errors that could be returned not just file doesn't
	// exist one
	_, err := os.Stat(path)
	return err
}

// createFile to create the UDR input files based on the testParams obj
func createFile(t *perftest.RatingParams) error {
	// check to see if the location exist, location specified must exist
	if err := exists(t.DropLocation); err != nil {
		return err
	}

	var filename string
	for i := 0; i < t.NumOfFiles; i++ {
		filename = t.DropLocation + "/" + t.FilenamePrefix + "-" + strconv.Itoa(i) + ".csv"

		fo, err := os.Create(filename)
		defer func() {
			if e := fo.Close(); e != nil {
				panic(e)
			}
		}()

		if err != nil {
			return err
		}

		for i := 0; i < t.NumRecordsPerFile; i++ {
			// No random, rate repeatly using the current timestamp for phase 1
			// 20060102150405 is const have to specify it this way, refer to
			// http://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
			tns := time.Now().Format("20060102150405.000")

			// replace the timestamp
			t.RawFields[t.TimpstampFieldIndex] = tns
			fo.WriteString(strings.Join(t.RawFields, ",") + "\n")
		}
	}

	return nil
}
