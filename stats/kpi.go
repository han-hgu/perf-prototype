package stats

import (
	"github.com/perf-prototype/perftest"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// UpdateDBParameters updates the database KPIs and configuration settings
// related to the performance
func (c *Controller) UpdateDBParameters(dbname string, dbp *perftest.DBParam) error {
	dbp.CompatibilityLevel = c.compatiblityLevel(dbname)
	return nil
}

func (c *Controller) compatiblityLevel(dbname string) (clevel uint8) {
	q := `SELECT compatibility_level FROM sys.databases WHERE name = '` + dbname + `'`
	c.getLastVal(q, []interface{}{&clevel})
	return clevel
}

// TrackKPI a container for any db related kpi tracking, this is db CPU, mem etc
// this is in the planning stage
func (c *Controller) TrackKPI(r perftest.Result) {
	// TODO this should use Brian's script for db perf figures
	mem, _ := mem.VirtualMemory()
	cpu, _ := cpu.Percent(0, false)
	r.AddDBServerMem(float32(mem.UsedPercent))
	r.AddDBServerCPU(float32(cpu[0]))
}
