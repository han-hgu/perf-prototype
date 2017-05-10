package perftest

type iController interface {
	UpdateRatingResult(ti *TestInfo, dbIDTracker *DBIDTracker) error
	UpdateBillingResult(ti *TestInfo, dbIDTracker *DBIDTracker) error
	UpdateBaselineIDs(t *DBIDTracker) error
	UpdateDBParameters(dbname string, dbp *DBParam) error
}
