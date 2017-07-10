package perftest

import "gopkg.in/mgo.v2/bson"

type storage interface {
	Initialize()
	Teardown()
	add(r Result) error
	get(id bson.ObjectId) (Result, error)
	getTestResultSVByTags(tags []string) ([]TestResultSV, error)
}
