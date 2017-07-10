package perftest

import (
	"errors"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

type mockStore struct {
	sync.RWMutex
	info map[bson.ObjectId]Result
}

func (s *mockStore) Initialize() {
	s.Lock()
	s.info = make(map[bson.ObjectId]Result, 0)
	s.Unlock()
}

func (s *mockStore) Teardown() {
	s.Lock()
	defer s.Unlock()
	s.info = make(map[bson.ObjectId]Result, 0)
}

func (s *mockStore) add(r Result) error {
	if s.info == nil {
		s.Initialize()
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.info[r.TestID()]; !ok {
		s.info[r.TestID()] = r
	}

	return nil
}

func (s *mockStore) get(id bson.ObjectId) (Result, error) {
	if s.info == nil {
		s.Initialize()
	}

	s.RLock()
	defer s.RUnlock()
	if _, ok := s.info[id]; !ok {
		return nil, errors.New("test doesn't exist")
	}

	return s.info[id], nil
}

func (s *mockStore) getTestResultSVByTags(tags []string) ([]TestResultSV, error) {
	if s.info == nil {
		return nil, nil
	}

	ret := make([]TestResultSV, 0)
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.info {
		trsv := TestResultSV{}

		trsv.ID = v.TestID()
		trsv.Md = v.MetaData()
		ret = append(ret, trsv)
	}

	return ret, nil
}
