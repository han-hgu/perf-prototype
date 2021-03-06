package perftest

import "sync"

type iController interface {
	UpdateRatingResult(*TestInfo, *DBIDTracker) error
	UpdateBillingResult(*TestInfo, *DBIDTracker) error
	UpdateBaselineIDs(*DBIDTracker) error
	UpdateDBParameters(*DBConf, *DBParams) error
	TrackKPI(wg *sync.WaitGroup, dbname string, relativeCPU *float32, logicalReads *uint64, logicalWrites *uint64, physicalReads *uint64)
	GetEIPOptions(ac *AppConf) error
	SnapshotIndexUsageStats(table string) error
	IndexUsageStatsComparison(tableBefore, tableAfter string) (map[string]map[string]interface{}, error)
}
