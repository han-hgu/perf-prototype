package perftest

type iController interface {
	UpdateRatingResult(ti *TestInfo, dbIDTracker *DBIDTracker) error
	UpdateBaselineIDs(t *DBIDTracker) error
}
