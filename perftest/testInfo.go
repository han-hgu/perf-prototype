package perftest

import (
	"sync"
	"time"
)

// Params interface to abstract out params
type Params interface {
	GetInfo() map[string]string
	GetTestID() string
	GetController() iController
	GetKeywords() map[string]string
	GetCollectionInterval() time.Duration
	GetDBConfig() *DBConf
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
	TestID             string            `json:"-"`
	DBConf             DBConf            `json:"db_config"`
	AdditionalInfo     map[string]string `json:"additional_info"`
	Keywords           map[string]string `json:"keywords"`
	DbController       iController       `json:"-"`
	CollectionInterval time.Duration     `json:"collection_interval"`
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

// RatingParams to hold the rating test parameters
type RatingParams struct {
	TestParams
	AmtFieldIndex       int      `json:"amount_field_index"`
	TimpstampFieldIndex int      `json:"timestamp_field_index"`
	NumOfFiles          uint32   `json:"number_of_files"`
	NumRecordsPerFile   int      `json:"number_of_records_per_file"`
	RawFields           []string `json:"raw_fields"`
	UseExistingFile     bool     `json:"use_existing_file"`
	DropLocation        string   `json:"drop_location"`
	FilenamePrefix      string   `json:"filename_prefix"`

	DBIDTracker *DBIDTracker `json:"-"`
}

// GetInfo to integrate RatingParams to Params interface
func (tp *TestParams) GetInfo() map[string]string {
	return tp.AdditionalInfo
}

// GetTestID returns the test ID
func (tp *TestParams) GetTestID() string {
	return tp.TestID
}

// GetController get the iController from the params
func (tp *TestParams) GetController() iController {
	return tp.DbController
}

// GetKeywords gets the keywords from the params
func (tp *TestParams) GetKeywords() map[string]string {
	return tp.Keywords
}

// GetCollectionInterval gets the collection interval from the params
func (tp *TestParams) GetCollectionInterval() time.Duration {
	return tp.CollectionInterval
}

// GetDBconfig gets the database configuration
func (tp *TestParams) GetDBConfig() *DBConf {
	return &(tp.DBConf)
}

// BillingParams to hold the billing test parameters
type BillingParams struct {
	TestParams
	OwnerName string `json:"owner_name"`
}

// DBParam stores the db parameters which could impact performance
type DBParam struct {
	CompatibilityLevel uint8 `json:"compatibility_level"`
}

// TestResult to store generic results
type TestResult struct {
	StartTime      time.Time         `json:"start_date"`
	Duration       string            `json:"test_duration,omitempty"`
	Done           bool              `json:"test_completed"`
	AdditionalInfo map[string]string `json:"additional_info"`
	Keywords       map[string]string `json:"keywords,omitempty"`
	CPUMax         float64           `json:"cpu_max(%)"`
	MemMax         float64           `json:"mem_max(%)"`
	DBParam        DBParam           `json:"database_parameters"`
}

// GetResult to integrate RatingResult to Result interface
func (rr *TestResult) GetResult() *TestResult {
	return rr
}

// RatingResult to save the rate information
type RatingResult struct {
	TestResult
	Rates                 []float32 `json:"udr_rates,omitempty"`
	MinRate               float32   `json:"-"`
	AvgRate               float32   `json:"udr_rate_avg,omitempty"`
	UDRProcessed          uint64    `json:"udr_created"`
	UDRExceptionProcessed uint64    `json:"udr_exception_created"`
	FilesCompleted        uint32    `json:"files_completed"`
}

// BillingResult to save the billing information
type BillingResult struct {
	TestResult
	UserPackageBilled          uint64    `json:"user_package_billed,omitempty"`
	InvoiceRenderDuration      string    `json:"invoice_render_duration,omitempty"`
	InvoiceRenderStartTime     time.Time `json:"-"`
	InvoiceRenderStartTimeOnce sync.Once `json:"-"`
	InvoiceRenderEndTime       time.Time `json:"-"`
	InvoiceRenderEndTimeOnce   sync.Once `json:"-"`
	BillingDuration            string    `json:"billing_duration,omitempty"`
	BillingStartTime           time.Time `json:"-"`
	BillingStartTimeOnce       sync.Once `json:"-"`
	BillingEndTime             time.Time `json:"-"`
	BillingEndTimeOnce         sync.Once `json:"-"`
	BillrunEndTime             time.Time `json:"-"`
	BillrunEndOnce             sync.Once `json:"-"`
	UserPackageBillRate        []uint32  `json:"user_package_bill_rate,omitempty"`
}
