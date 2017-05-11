package perftest

type iController interface {
	UpdateRatingResult(*TestInfo, *DBIDTracker) error
	UpdateBillingResult(*TestInfo, *DBIDTracker) error
	UpdateBaselineIDs(*DBIDTracker) error
	UpdateDBParameters(string, *DBParam) error
	TrackKPI(Result)
}
