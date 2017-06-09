package perftest

import (
	"errors"
	"log"
	"sync"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// MongoDBHost for mongo database host info
	MongoDBHost = "192.168.1.49"
	// DBName to save the perf results
	DBName = "EngageIPPerfResults"
	// CollectionName for the collection name for perf results
	CollectionName = "TestInfo"
)

// store to save test information
type store struct {
	mgoSession *mgo.Session

	metaDataInfoLock sync.RWMutex
	metaDataInfo     map[bson.ObjectId]Metadata

	sync.RWMutex
	info map[string]*TestInfo
}

// laod all meta data into memory
func (s *store) initialize() {
	db, err := mgo.Dial(MongoDBHost)
	db.SetMode(mgo.Monotonic, true)
	s.mgoSession = db
	if err != nil {
		log.Fatal("ERR: cannot connect to backbone database", err)
	}

	// load the meta data to the map, don't use cache since it has to be
	// completely loaded
	s.metaDataInfoLock.Lock()
	defer s.metaDataInfoLock.Unlock()
	currSession := s.mgoSession.Copy()
	defer currSession.Close()
}

func (s *store) teardown() {
	s.mgoSession.Close()
}

func (s *store) add(uuid string, t *TestInfo) error {
	if s.info == nil {
		s.info = make(map[string]*TestInfo)
	}

	s.Lock()
	defer s.Unlock()
	if _, ok := s.info[uuid]; !ok {
		s.info[uuid] = t
		return nil
	}

	return errors.New("test already exists")
}

// get testInfo from the store
func (s *store) get(uuid string) (TestInfo, error) {
	if s.info == nil {
		return TestInfo{}, errors.New("test doesn't exist")
	}

	s.RLock()
	defer s.RUnlock()
	if _, ok := s.info[uuid]; !ok {
		return TestInfo{}, errors.New("test doesn't exist")
	}

	return *s.info[uuid], nil
}

func (s *store) getAll() []map[string]interface{} {
	retVal := make([]map[string]interface{}, 0)
	if s.info == nil {
		return nil
	}

	s.RLock()
	defer s.RUnlock()
	for _, ti := range s.info {
		currTest := make(map[string]interface{})
		currTest["id"] = ti.Result.TestID()
		currTest["meta_data"] = ti.Result.MetaData()
		retVal = append(retVal, currTest)
	}

	return retVal
}

func (s *store) update(uuid string, t *TestInfo) error {
	if s.info == nil {
		return errors.New("update non-existing test result")
	}
	s.Lock()
	defer s.Unlock()

	if _, ok := s.info[uuid]; !ok {
		log.Println("WARNING: updating an non-existing test")
		return errors.New("update non-existing test result")
	}

	s.info[uuid] = t
	return nil
}
