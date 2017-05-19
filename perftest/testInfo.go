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
	CollectionInterval() string
	DBConfig() *DBConf
	AppConfig() *AppConf
}

// Result interface to abstract out results
type Result interface {
	Result() *TestResult
	AddAppServerCPU(float32)
	AddAppServerMem(float32)
	AddDBServerCPU(float32)
	AddDBServerMem(float32)
	AppServerStats() genericStats
	DBServerStats() genericStats
	TestID() string
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

// perfMonAppStats to hold the stats from perfmon
type perfMonAppStats struct {
	Mem float32 `json:"mem"`
	CPU float32 `json:"cpu"`
}

// AppConf for perfmon url
type AppConf struct {
	URL string `json:"url"`
}

// TestParams to hold common test parameters for all test types
type TestParams struct {
	ID             string            `json:"-"`
	DBConf         DBConf            `json:"db_config"`
	AppConf        AppConf           `json:"app_config"`
	AdditionalInfo map[string]string `json:"additional_info"`
	Kwords         map[string]string `json:"keywords"`
	DbController   iController       `json:"-"`
	CInterval      string            `json:"collection_interval"`
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
func (tp *TestParams) CollectionInterval() string {
	return tp.CInterval
}

// DBConfig gets the database configuration
func (tp *TestParams) DBConfig() *DBConf {
	return &(tp.DBConf)
}

// AppConfig gets the app server information
func (tp *TestParams) AppConfig() *AppConf {
	return &(tp.AppConf)
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

// genericStats for CPU and memory
type genericStats struct {
	CPU       []float32 `json:"cpu(%)"`
	CPUMaxium float32   `json:"cpu_max(%)"`
	Mem       []float32 `json:"mem(%)"`
	MemMaxium float32   `json:"mem_max(%)"`
}

// TestResult to store generic results
type TestResult struct {
	ID             string            `json:"ID"`
	StartTime      time.Time         `json:"start_date"`
	Duration       string            `json:"test_duration,omitempty"`
	Done           bool              `json:"test_completed"`
	AdditionalInfo map[string]string `json:"additional_info"`
	Keywords       map[string]string `json:"keywords,omitempty"`
	AppStats       genericStats      `json:"app_server_stats"`
	DBStats        genericStats      `json:"database_server_stats"`
	DBParam        DBParam           `json:"database_parameters"`
}

// TestID to get the test ID
func (rr *TestResult) TestID() string {
	return rr.ID
}

// Result to integrate RatingResult to Result interface
func (rr *TestResult) Result() *TestResult {
	return rr
}

// AddAppServerCPU adds a cpu sample for the app server
func (rr *TestResult) AddAppServerCPU(v float32) {
	if rr.AppStats.CPU == nil {
		rr.AppStats.CPU = make([]float32, 0)
	}

	rr.AppStats.CPU = append(rr.AppStats.CPU, v)
	if rr.AppStats.CPUMaxium < v {
		rr.AppStats.CPUMaxium = v
	}
}

// AddAppServerMem adds a memory sample for the app server
func (rr *TestResult) AddAppServerMem(v float32) {
	if rr.AppStats.Mem == nil {
		rr.AppStats.Mem = make([]float32, 0)
	}

	rr.AppStats.Mem = append(rr.AppStats.Mem, v)
	if rr.AppStats.MemMaxium < v {
		rr.AppStats.MemMaxium = v
	}
}

// AppServerStats to return the stats object for app server
func (rr *TestResult) AppServerStats() genericStats {
	return rr.AppStats
}

// AddDBServerCPU adds a cpu sample for the database
func (rr *TestResult) AddDBServerCPU(v float32) {
	if rr.DBStats.CPU == nil {
		rr.DBStats.CPU = make([]float32, 0)
	}

	rr.DBStats.CPU = append(rr.DBStats.CPU, v)
	if rr.DBStats.CPUMaxium < v {
		rr.DBStats.CPUMaxium = v
	}
}

// AddDBServerMem adds a memory sample for the database
func (rr *TestResult) AddDBServerMem(v float32) {
	if rr.DBStats.Mem == nil {
		rr.DBStats.Mem = make([]float32, 0)
	}

	rr.DBStats.Mem = append(rr.DBStats.Mem, v)
	if rr.DBStats.MemMaxium < v {
		rr.DBStats.MemMaxium = v
	}
}

// DBServerStats to return the stats object for db server
func (rr *TestResult) DBServerStats() genericStats {
	return rr.DBStats
}

// RatingResult to save the rate information
type RatingResult struct {
	TestResult
	Rates                 []float32 `json:"udr_rates,omitempty"`
	MinRate               float32   `json:"-"`
	AvgRate               float32   `json:"udr_rate_avg,omitempty"`
	UDRProcessedTrend     []uint64  `json:"udr_created_trend"`
	UDRProcessed          uint64    `json:"udr_created"`
	UDRExceptionProcessed uint64    `json:"udr_exception_created"`
	FilesCompleted        uint32    `json:"files_completed"`
}

// BillingResult to save the billing information
type BillingResult struct {
	TestResult
	UserPackageBilled        uint64    `json:"user_package_billed,omitempty"`
	InvoiceRenderDuration    string    `json:"invoice_render_duration,omitempty"`
	InvoiceRenderStartTime   time.Time `json:"-"`
	InvoiceRenderEndTime     time.Time `json:"-"`
	InvoiceRenderEndTimeOnce sync.Once `json:"-"`
	BillingDuration          string    `json:"billing_duration,omitempty"`
	BillingStartTime         time.Time `json:"-"`
	BillingEndTime           time.Time `json:"-"`
	BillingEndTimeOnce       sync.Once `json:"-"`
	BillrunEndTime           time.Time `json:"-"`
	BillrunEndOnce           sync.Once `json:"-"`
	UserPackageBillRate      []uint32  `json:"user_package_bill_rate,omitempty"`
}
