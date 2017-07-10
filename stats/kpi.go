package stats

import (
	"log"
	"sync"

	"github.com/perf-prototype/perftest"
)

// UpdateDBParameters updates the database KPIs and configuration settings
// related to the performance
func (c *Controller) UpdateDBParameters(dbc *perftest.DBConf, dbp *perftest.DBParams) error {
	dbp.Database = dbc.Database
	dbp.AddtionalInfo = dbc.AddtionalInfo
	dbp.URL = dbc.URL
	dbp.CompatibilityLevel = c.compatiblityLevel(dbp.Database)
	return nil
}

func (c *Controller) compatiblityLevel(dbname string) (clevel uint8) {
	q := `SELECT compatibility_level FROM sys.databases WHERE name = '` + dbname + `'`
	c.getLastVal(q, []interface{}{&clevel})
	return clevel
}

// TrackKPI a container for any db related kpi tracking, this is db CPU, mem etc
// this is in the planning stage
func (c *Controller) TrackKPI(wg *sync.WaitGroup, dbname string, cpu *float32, logicalReads *uint64, logicalWrites *uint64, physicalReads *uint64) {
	if wg != nil {
		defer wg.Done()
	}

	q := `SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
	      ;WITH DB_CPU_Stats
		  AS
		  (
		    SELECT DatabaseID, isnull(DB_Name(DatabaseID),
				   case DatabaseID when 32767 then 'Internal ResourceDB' else CONVERT(varchar(255),DatabaseID)end) AS [DatabaseName],
		           SUM(total_worker_time) AS [CPU_Time_Ms], SUM(total_logical_reads)  AS [Logical_Reads],
		           SUM(total_logical_writes)  AS [Logical_Writes], SUM(total_logical_reads+total_logical_writes)  AS [Logical_IO],
		           SUM(total_physical_reads)  AS [Physical_Reads], SUM(total_elapsed_time)  AS [Duration_MicroSec],
		           SUM(total_clr_time)  AS [CLR_Time_MicroSec], SUM(total_rows)  AS [Rows_Returned],
		           SUM(execution_count)  AS [Execution_Count], count(*) 'Plan_Count'
		    FROM sys.dm_exec_query_stats AS qs
		    CROSS APPLY (
				   SELECT CONVERT(int, value) AS [DatabaseID]
				   FROM sys.dm_exec_plan_attributes(qs.plan_handle)
				   WHERE attribute = N'dbid') AS F_DB
				   GROUP BY DatabaseID
		  )
		  SELECT DatabaseName,
		  CAST([CPU_Time_Ms] * 1.0 / SUM(case [CPU_Time_Ms] when 0 then 1 else [CPU_Time_Ms] end) OVER() * 100.0 AS DECIMAL(5, 2)) AS [CPU_Percent],
		  [Logical_Reads],
		  [Physical_Reads],
		  [Logical_Writes]
		  FROM DB_CPU_Stats;`

	rows, err := c.db.Query(q)
	if err != nil {
		log.Fatalf("ERR: Database KPI query returns error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var dbc string
		rowErr := rows.Scan(&dbc, cpu, logicalReads, physicalReads, logicalWrites)
		if rowErr != nil {
			log.Printf("WARNING: get error scanning the rows for db KPI query: %v\n", rowErr)
			continue
		}

		if dbc == dbname {
			return
		}
	}

	log.Printf("WARNING: database %v not found by KPI query\n", dbname)
}
