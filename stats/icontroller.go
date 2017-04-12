package stats

// Params to store parameters for stats collection, every perftest.worker holds
// one of these
// NOTE: no protection for mutithreading, each field should have one thread
// responsible for it
type Params struct {
	FilenamePrefix     string
	LastEventLogID     uint64
	LastUdrID          uint64
	LastUdrExceptionID uint64
}

type statsController interface {
	TearDown()
	GetRate(*Params)
}
