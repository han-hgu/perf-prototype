package perftest

import (
	"encoding/json"
	"log"

	"github.com/allegro/bigcache"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// MongoDBHost for mongo database host info
	MongoDBHost = "192.168.1.49"
	// DBName to save the perf results
	DBName = "EIPPerfFramework"
	// CollectionName for the collection name for perf results
	CollectionName = "TestResult"
)

// Store to save test information
type Store struct {
	mgoSession *mgo.Session
	cache      *bigcache.BigCache
}

// Initialize initializes the store
func (s *Store) Initialize() {
	db, err := mgo.Dial(MongoDBHost)
	if err != nil {
		log.Fatalf("ERR: cannot connect to backbone database: %v", err)
	}
	db.SetMode(mgo.Monotonic, true)
	db.SetSafe(&mgo.Safe{})
	s.mgoSession = db

	// should be able to evict immediately
	s.cache, err = bigcache.NewBigCache(bigcache.DefaultConfig(0))
	if err != nil {
		log.Fatalf("ERR: fail to intialize cache: %v", err)
	}
}

// Teardown shuts down the store
func (s *Store) Teardown() {
	s.mgoSession.Close()
}

func (s *Store) cacheResult(r Result) error {
	// cache it
	b, e := json.Marshal(r)
	if e != nil {
		return e
	}
	s.cache.Set(r.TestID().Hex(), b)
	return nil
}

func (s *Store) add(r Result) error {
	// cache it
	if e := s.cacheResult(r); e != nil {
		log.Printf("WARNING: failed to cache result, %v", e)
	}

	// store it
	session := s.mgoSession.Copy()
	defer session.Close()

	collection := session.DB(DBName).C(CollectionName)
	err := collection.Insert(r)
	if mgo.IsDup(err) {
		return nil
	}

	return err
}

// get testInfo from the store
func (s *Store) get(id bson.ObjectId) (r Result, err error) {
	// check if it is already cached, if there is any error grab from
	// database
	if rb, e := s.cache.Get(id.Hex()); e == nil {
		var gr GeneralResult
		if e := json.Unmarshal(rb, &gr); e == nil {
			if gr.MetaData().Type == "rating" {
				var rr RatingResult
				if e := json.Unmarshal(rb, &rr); e == nil {
					return &rr, nil
				}
			} else if gr.MetaData().Type == "billing" {
				var br BillingResult
				if e := json.Unmarshal(rb, &br); e == nil {
					return &br, nil
				}
			}
		}
	}

	session := s.mgoSession.Copy()
	defer session.Close()

	rt := TestResult{}
	collection := session.DB(DBName).C(CollectionName)
	err = collection.FindId(id).One(&rt)

	if rt.Type == "rating" {
		r = new(RatingResult)
	} else if rt.Type == "billing" {
		r = new(BillingResult)
	} else {
		r = new(TestResult)
	}

	err = collection.FindId(id).One(r)

	// cache result
	if e := s.cacheResult(r); e != nil {
		log.Printf("WARNING: failed to cache result, %v", e)
	}
	return r, err
}

// getTestResultSVByTags returns the results matching all specified tags, each
// returned element contains a short version of test result
func (s *Store) getTestResultSVByTags(tags []string) ([]TestResultSV, error) {
	r := make([]TestResultSV, 0)
	session := s.mgoSession.Copy()
	defer session.Close()

	collection := session.DB(DBName).C(CollectionName)
	var findQ bson.M
	if len(tags) != 0 {
		findQ = bson.M{"meta_data.tags": bson.M{"$all": tags}}
	}

	selectQ := bson.M{"meta_data": 1}
	err := collection.Find(findQ).Select(selectQ).All(&r)
	return r, err
}
