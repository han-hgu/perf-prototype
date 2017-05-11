package stats

import (
	"fmt"

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
	mem, _ := mem.VirtualMemory()
	cpu, _ := cpu.Percent(0, false)
	if mem.UsedPercent > r.MemMax() {
		r.SetMemMax(mem.UsedPercent)
	}

	if cpu[0] > r.CPUMax() {
		r.SetCPUMax(cpu[0])
	}

	fmt.Printf("MemUsedPercent:\t%f%%\n", mem.UsedPercent)
	fmt.Printf("CPUPercent:\t%f%%\n", cpu[0])
}
