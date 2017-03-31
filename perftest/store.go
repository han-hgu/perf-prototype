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

func (s *store) add(uuid string, t *TestInfo) {
	s.Lock()
	defer s.Unlock()
	s.info[uuid] = t
}

// get testInfo from the store
func (s *store) get(uuid string) (TestInfo, error) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.info[uuid]; !ok {
		return TestInfo{}, errors.New("test doesn't exist")
	}

	return *s.info[uuid], nil
}

func (s *store) update(uuid string, t TestInfo) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.info[uuid]; !ok {
		log.Println("WARNING: updating an non-existing test")
		return errors.New("update non-existing test result")
	}

	s.info[uuid] = &t
	return nil
}
