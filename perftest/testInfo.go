package perftest

import "time"

// Params interface to abstract out params
type Params interface {
	GetInfo() map[string]string
	GetTestID() string
	GetController() iController
}

// Result interface to abstract out results
type Result interface {
	GetResult() *TestResult
}

// TestInfo stores all the test related information
type TestInfo struct {
	Params Params
	Result Result
}

// DBConf for database connection information
type DBConf struct {
	Server   string `json:"ip"`
	Port     int    `json:"port"`
	Database string `json:"db_name"`
	UID      string `json:"uid"`
	Pwd      string `json:"password"`
}

// TestParams to hold common test parameters for all test types
type TestParams struct {
	TestID         string            `json:"-"`
	DBConf         DBConf            `json:"db_config"`
	AdditionalInfo map[string]string `json:"additional_info"`
	DbController   iController       `json:"-"`
}

// DBIDTracker keeps track of the last database table IDs examined
type DBIDTracker struct {
	EventLogLastProcessed     uint64
	EventLogCurrent           uint64
	EventlogStarted           uint64
	UDRLastProcessed          uint64
	UDRCurrent                uint64
	UDRStarted                uint64
	UDRExceptionLastProcessed uint64
	UDRExceptionCurrent       uint64
	UDRExceptionStarted       uint64
	TimePrevious              time.Time
}

// RatingParams to hold the testing parameters
type RatingParams struct {
	TestParams
	AmtFieldIndex       int           `json:"amount_field_index"`
	TimpstampFieldIndex int           `json:"timestamp_field_index"`
	NumOfFiles          uint32        `json:"number_of_files"`
	NumRecordsPerFile   int           `json:"number_of_records_per_file"`
	RawFields           []string      `json:"raw_fields"`
	UseExistingFile     bool          `json:"use_existing_file"`
	DropLocation        string        `json:"drop_location"`
	FilenamePrefix      string        `json:"filename_prefix"`
	CollectionInterval  time.Duration `json:"collection_interval"`
	DBIDTracker         *DBIDTracker  `json:"-"`
}

// GetInfo to integrate RatingParams to Params interface
func (rp *RatingParams) GetInfo() map[string]string {
	return rp.AdditionalInfo
}

// GetTestID returns the test ID
func (rp *RatingParams) GetTestID() string {
	return rp.TestID
}

// GetController get the iController from the params
func (rp *RatingParams) GetController() iController {
	return rp.DbController
}

// TestResult to store generic results
type TestResult struct {
	StartTime      time.Time         `json:"-"`
	Done           bool              `json:"test_completed"`
	AdditionalInfo map[string]string `json:"additional_info"`
}

// RatingResult to save the rate information
type RatingResult struct {
	TestResult
	Rates                 []float32 `json:"rates"`
	MinRate               float32   `json:"MIN_rate"`
	AvgRate               float32   `json:"AVG_rate"`
	UDRProcessed          uint64    `json:"udr_created"`
	UDRExceptionProcessed uint64    `json:"udr_exception_created"`
	FilesCompleted        uint32    `json:"files_completed"`
}

// GetResult to integrate RatingResult to Result interface
func (rr *RatingResult) GetResult() *TestResult {
	return &rr.TestResult
}
