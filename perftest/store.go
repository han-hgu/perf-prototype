package perftest

import (
	"errors"
	"log"
	"sync"
)

// store to save test information
type store struct {
	sync.RWMutex
	info map[string]*TestInfo
}

// laod all meta data into memory
func (s *store) initialize() {

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

func (s *store) getAll() []Metadata {
	r := make([]Metadata, 0)
	if s.info == nil {
		return nil
	}

	s.RLock()
	defer s.RUnlock()
	for _, ti := range s.info {
		r = append(r, ti.Result.MetaData())
	}

	return r
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
