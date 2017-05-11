package perftest

import (
	"sync"
	"time"
)

// Params interface to abstract out params
type Params interface {
	Info() map[string]string
	TestID() string
	Controller() iController
	Keywords() map[string]string
	CollectionInterval() time.Duration
	DBConfig() *DBConf
}

// Result interface to abstract out results
type Result interface {
	Result() *TestResult
	CPUMax() float64
	MemMax() float64
	SetCPUMax(float64)
	SetMemMax(float64)
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
	ID             string            `json:"-"`
	DBConf         DBConf            `json:"db_config"`
	AdditionalInfo map[string]string `json:"additional_info"`
	Kwords         map[string]string `json:"keywords"`
	DbController   iController       `json:"-"`
	CInterval      time.Duration     `json:"collection_interval"`
}

// Info returns the AdditionalInfo field
func (tp *TestParams) Info() map[string]string {
	return tp.AdditionalInfo
}

// TestID returns the test ID
func (tp *TestParams) TestID() string {
	return tp.ID
}

// Controller returns the iController from the params
func (tp *TestParams) Controller() iController {
	return tp.DbController
}

// Keywords returns the keywords from the params
func (tp *TestParams) Keywords() map[string]string {
	return tp.Kwords
}

// CollectionInterval gets the collection interval from the params
func (tp *TestParams) CollectionInterval() time.Duration {
	return tp.CInterval
}

// DBConfig gets the database configuration
func (tp *TestParams) DBConfig() *DBConf {
	return &(tp.DBConf)
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

// RatingParams holds rating test parameters
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

// BillingParams holds billing test parameters
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
	CPUMaxium      float64           `json:"cpu_max(%)"`
	MemMaxium      float64           `json:"mem_max(%)"`
	DBParam        DBParam           `json:"database_parameters"`
}

// Result to integrate RatingResult to Result interface
func (rr *TestResult) Result() *TestResult {
	return rr
}

// CPUMax returns rr.CPUMax
func (rr *TestResult) CPUMax() float64 {
	return rr.CPUMaxium
}

// MemMax returns rr.MemMax
func (rr *TestResult) MemMax() float64 {
	return rr.MemMaxium
}

// SetCPUMax sets the CPUMax field
func (rr *TestResult) SetCPUMax(v float64) {
	rr.CPUMaxium = v
}

// SetMemMax sets the MemMax field
func (rr *TestResult) SetMemMax(v float64) {
	rr.MemMaxium = v
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
