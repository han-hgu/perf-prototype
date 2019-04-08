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

// SnapshotIndexUsageStats takes a snapshot of the current system for index usage and store it to table
func (c *Controller) SnapshotIndexUsageStats(tableName string) error {
	q := `select OBJECT_NAME(object_id) 'table', database_id, object_id, index_id, user_seeks, user_scans, user_lookups, user_updates, last_user_seek, last_user_scan, last_user_lookup, last_user_update,getDate() 'DateGenerated'
		into ` + tableName +
		` from sys.dm_db_index_usage_stats
		where database_id = DB_ID()
		order by user_seeks desc`

	if _, err := c.db.Exec(q); err != nil {
		return err
	}

	return nil
}

func (c *Controller) IndexUsageStatsComparison(tableBefore, tableAfter string) (map[string]map[string]interface{}, error) {
	ret := make(map[string]map[string]interface{}, 0)

	var queryParams []interface{}
	var table string
	queryParams = append(queryParams, &table)

	var index_id uint64
	queryParams = append(queryParams, &index_id)

	var user_seek_diff uint64
	queryParams = append(queryParams, &user_seek_diff)

	var user_scan_diff uint64
	queryParams = append(queryParams, &user_scan_diff)

	var user_lookup_diff uint64
	queryParams = append(queryParams, &user_lookup_diff)

	var user_index_read_diff uint64
	queryParams = append(queryParams, &user_index_read_diff)

	var user_update_diff uint64
	queryParams = append(queryParams, &user_update_diff)

	var table_read_diff uint64
	queryParams = append(queryParams, &table_read_diff)

	q := `;with t1 as (select
	        poidx.[table], poidx.index_id
	            ,isnull(poidx.user_seeks, 0) - isnull(preidx.user_seeks, 0) user_seek_diff
	            ,isnull(poidx.user_scans, 0) - isnull(preidx.user_scans, 0) user_scan_diff
	            ,isnull(poidx.user_lookups, 0) - isnull(preidx.user_lookups, 0) user_lookup_diff
	            ,isnull(poidx.user_seeks, 0) - isnull(preidx.user_seeks, 0)
	            + isnull(poidx.user_scans, 0) - isnull(preidx.user_scans, 0)
	            + isnull(poidx.user_lookups, 0) - isnull(preidx.user_lookups, 0) user_index_read_diff
	            ,isnull(poidx.user_updates, 0) - isnull(preidx.user_updates, 0) user_update_diff
	            --,*
	        from ` + tableAfter + " poidx " +
		"full join " + tableBefore + " preidx " +
		`on poidx.object_id = preidx.object_id and poidx.index_id = preidx.index_id)
          select *, sum(user_index_read_diff) over (partition by [table]) table_read_diff from t1 where user_index_read_diff > 0`

	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		rowErr := rows.Scan(queryParams...)
		if rowErr != nil {
			return nil, rowErr
		}

		kpis := make(map[string]interface{}, 0)
		kpis["index_id"] = index_id
		kpis["user_seek_diff"] = user_seek_diff
		kpis["user_scan_diff"] = user_scan_diff
		kpis["user_lookup_diff"] = user_lookup_diff
		kpis["user_index_read_diff"] = user_index_read_diff
		kpis["user_update_diff"] = user_update_diff
		kpis["table_read_diff"] = table_read_diff
		ret[table] = kpis
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// drop the tracking tables
	dropTablesQ := `drop table ` + tableBefore + `; drop table ` + tableAfter + `;`
	if _, err := c.db.Exec(dropTablesQ); err != nil {
		return nil, err
	}

	return ret, nil
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
