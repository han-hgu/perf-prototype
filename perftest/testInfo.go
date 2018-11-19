package perftest

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Params interface to abstract out params
type Params interface {
	Comment() string
	TestID() bson.ObjectId
	Controller() iController
	Keywords() []string
	CollectionInterval() string
	DBConfig() *DBConf
	AppConfig() *AppConf
	ChartTitle() string
}

// Result interface to abstract out results
type Result interface {
	Result() *TestResult
	AddAppServerCPU(float32)
	AddDBCPU(float32)
	AddAppServerMem(float32)
	AddLogicalReads(v uint64)
	AddLogicalWrites(v uint64)
	AddPhysicalReads(v uint64)
	AppServerStats() *GenericStats
	DBServerStats() *DBStats
	TestID() bson.ObjectId
	CollectionInterval() string
	ChartTitle() string
	MetaData() Metadata
}

// TestInfo stores all the test related information
type TestInfo struct {
	Params Params
	Result Result
}

// PerfMonStats stores the stats from perfmon
type PerfMonStats struct {
	Mem float32 `json:"mem"`
	CPU float32 `json:"cpu"`
}

// DBConf for db connection information, and additional info later saved in
// result metadata portion
type DBConf struct {
	Server        string            `json:"ip"`
	Port          int               `json:"port"`
	Database      string            `json:"db_name"`
	UID           string            `json:"uid"`
	Pwd           string            `json:"password"`
	URL           string            `json:"perfmon_url"`
	AddtionalInfo map[string]string `json:"additional_info"`
}

// ChartConf for the configuration of how the data would be plotted
type ChartConf struct {
	Title string `json:"title"`
}

// AppConf for perfmon url
type AppConf struct {
	Version       string                 `json:"version" bson:"version"`
	Options       map[string]string      `json:"EIP_option" bson:"EIP_option"`
	URL           string                 `json:"perfmon_url" bson:"perfmon_url"`
	SysInfo       map[string]interface{} `json:"sys_info" bson:"sys_info"`
	AddtionalInfo map[string]string      `json:"additional_info,omitempty" bson:"additional_info,omitempty"`
}

// TestParams to hold common test parameters for all test types
type TestParams struct {
	ID           bson.ObjectId `json:"-"`
	DBConf       DBConf        `json:"db_config"`
	AppConf      AppConf       `json:"app_config"`
	ChartConf    ChartConf     `json:"chart_config"`
	Cmt          string        `json:"comment"`
	Kwords       []string      `json:"tags"`
	CInterval    string        `json:"collection_interval"`
	DbController iController   `json:"-"`
}

// Comment returns the Comment field
func (tp *TestParams) Comment() string {
	return tp.Cmt
}

// TestID returns the test ID
func (tp *TestParams) TestID() bson.ObjectId {
	return tp.ID
}

// Controller returns the iController from the params
func (tp *TestParams) Controller() iController {
	return tp.DbController
}

// Keywords returns the tags from the params
func (tp *TestParams) Keywords() []string {
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

// ChartTitle gets the title for chart
func (tp *TestParams) ChartTitle() string {
	return tp.ChartConf.Title
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
	NumOfUDRRecords     uint64   `json:"number_of_records"`
	NumRecordsPerFile   int      `json:"number_of_records_per_file"`
	RawFields           []string `json:"raw_fields"`
	UseExistingFile     bool     `json:"use_existing_file"`
	DropLocation        string   `json:"drop_location"`
	FilenamePrefix      string   `json:"filename_prefix"`
}

// BillingParams holds billing test parameters
type BillingParams struct {
	TestParams
	OwnerName string `json:"owner_name"`
}

// DBParams stores the db parameters
type DBParams struct {
	Database           string                 `json:"db_name" bson:"db_name"`
	URL                string                 `json:"perfmon_url" bson:"perfmon_url"`
	CompatibilityLevel uint8                  `json:"compatibility_level" bson:"compatibility_level"`
	SysInfo            map[string]interface{} `json:"sys_info" bson:"sys_info"`
	AddtionalInfo      map[string]string      `json:"additional_info,omitempty" bson:"additional_info,omitempty"`
}

// GenericStats for CPU and memory
type GenericStats struct {
	CPU       []float32 `json:"cpu(%),omitempty"`
	CPUMaxium float32   `json:"cpu_max(%)"`
	Mem       []float32 `json:"mem(%),omitempty"`
	MemMaxium float32   `json:"mem_max(%)"`
}

// DBStats holds database stats
type DBStats struct {
	GenericStats
	LReadsBase   uint64   `json:"-"`
	LReadsTotal  uint64   `json:"logical_reads_total"`
	LReads       []uint64 `json:"logical_reads,omitempty"`
	LWritesBase  uint64   `json:"-"`
	LWritesTotal uint64   `json:"logical_writes_total"`
	LWrites      []uint64 `json:"logical_writes,omitempty"`
	PReadsBase   uint64   `json:"-"`
	PReadsTotal  uint64   `json:"physical_reads_total"`
	PReads       []uint64 `json:"physical_reads,omitempty"`
}

// Metadata to store the metadata for the test, this is to make the search and
// display more user friendly
type Metadata struct {
	Type      string    `json:"test_type" bson:"test_type"`
	StartTime time.Time `json:"start_date" bson:"start_date"`
	Done      bool      `json:"test_completed" bson:"test_completed"`
	Duration  string    `json:"test_duration,omitempty" bson:"test_duration"`
	Keywords  []string  `json:"tags,omitempty" bson:"tags"`
	Cmt       string    `json:"comment,omitempty" bson:"comment"`
	CInterval string    `json:"collection_interval" bson:"collection_interval"`
	AppConf   AppConf   `json:"app_param" bson:"app_param"`
	DBParams  DBParams  `json:"db_param" bson:"db_param"`
	CTitle    string    `json:"chart_title" bson:"chart_title"`
}

// TestResultSV stands for short-version test result, only this will reside in
// the memory from a search, the complete result should be stored in a cache
// to control the memory footprint of the program
type TestResultSV struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
	Md Metadata      `json:"meta_data" bson:"meta_data"`
}

// TestResult to store generic results
type TestResult struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Metadata `json:"meta_data" bson:"meta_data"`
	AppStats GenericStats `json:"app_stats" bson:"app_stats"`
	DBStats  DBStats      `json:"db_stats" bson:"db_stats"`
}

// TestID to get the test ID
func (tr *TestResult) TestID() bson.ObjectId {
	return tr.ID
}

// CollectionInterval returns the collection interval
func (tr *TestResult) CollectionInterval() string {
	return tr.CInterval
}

// ChartTitle returns the title for the chart
func (tr *TestResult) ChartTitle() string {
	if tr.CTitle == "" {
		return tr.TestID().Hex()
	}

	return tr.CTitle
}

// Result to integrate RatingResult to Result interface
func (tr *TestResult) Result() *TestResult {
	return tr
}

// AddAppServerCPU adds a cpu sample for the app server
func (tr *TestResult) AddAppServerCPU(v float32) {
	if tr.AppStats.CPU == nil {
		tr.AppStats.CPU = make([]float32, 0)
	}

	tr.AppStats.CPU = append(tr.AppStats.CPU, v)
	if tr.AppStats.CPUMaxium < v {
		tr.AppStats.CPUMaxium = v
	}
}

// AddDBCPU adds a cpu sample for the database
func (tr *TestResult) AddDBCPU(v float32) {
	if tr.DBStats.CPU == nil {
		tr.DBStats.CPU = make([]float32, 0)
	}

	tr.DBStats.CPU = append(tr.DBStats.CPU, v)
	if tr.DBStats.CPUMaxium < v {
		tr.DBStats.CPUMaxium = v
	}
}

// AddAppServerMem adds a memory sample for the app server
func (tr *TestResult) AddAppServerMem(v float32) {
	if tr.AppStats.Mem == nil {
		tr.AppStats.Mem = make([]float32, 0)
	}

	tr.AppStats.Mem = append(tr.AppStats.Mem, v)
	if tr.AppStats.MemMaxium < v {
		tr.AppStats.MemMaxium = v
	}
}

// AddLogicalReads adds the number of logical reads per collection interval to
// test result
func (tr *TestResult) AddLogicalReads(v uint64) {
	if tr.DBStats.LReads == nil {
		tr.DBStats.LReads = make([]uint64, 0)
	}

	tr.DBStats.LReads = append(tr.DBStats.LReads, v)
	tr.DBStats.LReadsTotal += v
}

// AddLogicalWrites adds the number of logical per collection interval to test
// result
func (tr *TestResult) AddLogicalWrites(v uint64) {
	if tr.DBStats.LWrites == nil {
		tr.DBStats.LWrites = make([]uint64, 0)
	}

	tr.DBStats.LWrites = append(tr.DBStats.LWrites, v)
	tr.DBStats.LWritesTotal += v
}

// AddPhysicalReads adds the number of physical reads per collection to test
// result
func (tr *TestResult) AddPhysicalReads(v uint64) {
	if tr.DBStats.PReads == nil {
		tr.DBStats.PReads = make([]uint64, 0)
	}

	tr.DBStats.PReads = append(tr.DBStats.PReads, v)
	tr.DBStats.PReadsTotal += v
}

// MetaData returns the metadata section from the test result
func (tr *TestResult) MetaData() Metadata {
	return tr.Metadata
}

// FetchAppServerCPUStats fetches app server CPU stats
func FetchAppServerCPUStats(r Result) []float32 {
	return r.AppServerStats().CPU
}

// FetchAppServerMemStats fetches app server mem stats
func FetchAppServerMemStats(r Result) []float32 {
	return r.AppServerStats().Mem
}

// FetchDBServerCPUStats fetches db CPU stats
func FetchDBServerCPUStats(r Result) []float32 {
	return r.DBServerStats().CPU
}

// FetchDBServerMemStats fetches db mem stats
func FetchDBServerMemStats(r Result) []float32 {
	return r.DBServerStats().Mem
}

// FetchDBServerLReads fetches db logical reads
func FetchDBServerLReads(r Result) []uint64 {
	return r.DBServerStats().LReads
}

// FetchDBServerPReads fetches db physical reads
func FetchDBServerPReads(r Result) []uint64 {
	return r.DBServerStats().PReads
}

// FetchDBServerLWrites fetches db logical writes
func FetchDBServerLWrites(r Result) []uint64 {
	return r.DBServerStats().LWrites
}

// FetchRates fetches rates
func FetchRates(r Result) []float32 {
	rr, ok := r.(*RatingResult)
	if !ok {
		panic("ERR: Fetch rates from a non-rating result")
	}
	return rr.Rates
}

// FetchUDRProcessedTrend fetches udr processed trend
func FetchUDRProcessedTrend(r Result) []uint64 {
	rr, ok := r.(*RatingResult)
	if !ok {
		panic("ERR: Fetch rates from a non-rating result")
	}
	return rr.UDRProcessedTrend
}

// AppServerStats returns app server stats
func (tr *TestResult) AppServerStats() *GenericStats {
	return &(tr.AppStats)
}

// AddDBServerCPU adds db cpu sample
func (tr *TestResult) AddDBServerCPU(v float32) {
	if tr.DBStats.CPU == nil {
		tr.DBStats.CPU = make([]float32, 0)
	}

	tr.DBStats.CPU = append(tr.DBStats.CPU, v)
	if tr.DBStats.CPUMaxium < v {
		tr.DBStats.CPUMaxium = v
	}
}

// AddDBServerMem adds db memory sample
func (tr *TestResult) AddDBServerMem(v float32) {
	if tr.DBStats.Mem == nil {
		tr.DBStats.Mem = make([]float32, 0)
	}

	tr.DBStats.Mem = append(tr.DBStats.Mem, v)
	if tr.DBStats.MemMaxium < v {
		tr.DBStats.MemMaxium = v
	}
}

// DBServerStats returns db stats object
func (tr *TestResult) DBServerStats() *DBStats {
	return &(tr.DBStats)
}

// GeneralResult enables us to blindly unmarshal and check the test type
type GeneralResult struct {
	TestResult `bson:",inline"`
}

// RatingResult stores rating stats
type RatingResult struct {
	TestResult            `bson:",inline"`
	MinRate               float32   `json:"-"`
	AvgRate               float32   `json:"udr_rate_avg,omitempty" bson:"udr_rate_avg"`
	UDRProcessed          uint64    `json:"udr_created" bson:"udr_created"`
	UDRExceptionProcessed uint64    `json:"udr_exception_created" bson:"udr_exception_created"`
	FilesCompleted        uint32    `json:"files_completed" bson:"files_completed"`
	Rates                 []float32 `json:"udr_rates,omitempty" bson:"udr_rates"`
	UDRProcessedTrend     []uint64  `json:"udr_created_trend,omitempty" bson:"udr_created_trend"`
}

// BillingResult stores billing stats
type BillingResult struct {
	TestResult                 `bson:",inline"`
	OwnerName                  string    `json:"owner_name" bson:"owner_name"`
	BillingDuration            string    `json:"billing_duration,omitempty" bson:"billing_duration"`
	InvoiceRenderDuration      string    `json:"invoice_render_duration,omitempty" bson:"invoice_render_duration"`
	InvoiceRenderStartTime     time.Time `json:"-"`
	InvoiceRenderEndTime       time.Time `json:"-"`
	InvoiceRenderEndTimeOnce   sync.Once `json:"-" bson:"-"`
	BillingStartTime           time.Time `json:"-"`
	BillingEndTime             time.Time `json:"-"`
	BillingEndTimeOnce         sync.Once `json:"-" bson:"-"`
	BillrunEndTime             time.Time `json:"-"`
	BillrunEndOnce             sync.Once `json:"-" bson:"-"`
	UserPackagesBilled         []uint64  `json:"user_packages_billed,omitempty" bson:"user_packages_billed"`
	UserServicesBilled         []uint64  `json:"user_services_billed,omitempty" bson:"user_services_billed"`
	UsersBilled                []uint64  `json:"users_billed,omitempty" bson:"users_billed"`
	UsageInvoicersBilled       []uint64  `json:"usage_invoicers_billed,omitempty" bson:"usage_invoicers_billed"`
	InvoicesClosed             []uint64  `json:"invoices_closed,omitempty" bson:"invoices_closed"`
	UsageTranscationsGenerated []uint64  `json:"usage_transactions_generated,omitempty" bson:"usage_transactions_generated"`
	MRCTransactionsGenerated   []uint64  `json:"mrc_transactions_generated,omitempty" bson:"mrc_transactions_generated"`
	BillUDRActionCompleted     []uint64  `json:"bill_udr_actions_completed" bson:"bill_udr_actions_completed"`
}
