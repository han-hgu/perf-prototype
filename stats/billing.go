package stats

import (
	"fmt"
	"time"

	"github.com/perf-prototype/perftest"
)

// UpdateBillingResult to update the billing results
func (c *Controller) UpdateBillingResult(ti *perftest.TestInfo, dbIDTracker *perftest.DBIDTracker) error {
	tr := ti.Result.(*perftest.BillingResult)
	if tr.Done {
		return nil
	}

	tp := ti.Params.(*perftest.BillingParams)
	dbIDTracker.EventLogCurrent = c.getLastEventLogID()

	if c.BillingStarted(tp.OwnerName, dbIDTracker.EventLogLastProcessed) {
		tr.BillingStartTimeOnce.Do(func() {
			tr.BillingStartTime = c.BillingStartTime(tp.OwnerName, dbIDTracker.EventLogLastProcessed)
		})
	}

	if c.BillingFinished(tp.OwnerName, dbIDTracker.EventLogLastProcessed) {
		tr.BillingEndTimeOnce.Do(func() {
			tr.BillingEndTime = c.BillingEndTime(tp.OwnerName, dbIDTracker.EventLogLastProcessed)
			tr.BillingDuration = tr.BillingEndTime.Sub(tr.BillingStartTime).String()
		})
	}

	if c.InvoiceRenderingStarted(tp.OwnerName, dbIDTracker.EventlogStarted) {
		tr.InvoiceRenderStartTimeOnce.Do(func() {
			tr.InvoiceRenderStartTime = c.InvoiceRenderingStartTime(tp.OwnerName, dbIDTracker.EventLogLastProcessed)
		})
	}

	if c.InvoiceRenderingFinished(tp.OwnerName, dbIDTracker.EventlogStarted) {
		tr.InvoiceRenderEndTimeOnce.Do(func() {
			tr.InvoiceRenderEndTime = c.InvoiceRenderingEndTime(tp.OwnerName, dbIDTracker.EventLogLastProcessed)
			tr.InvoiceRenderDuration = tr.InvoiceRenderEndTime.Sub(tr.InvoiceRenderStartTime).String()
		})
	}

	if c.BillrunFinished(tp.OwnerName, dbIDTracker.EventLogLastProcessed) {
		tr.BillrunEndOnce.Do(func() {
			tr.BillrunEndTime = c.BillrunEndTime(tp.OwnerName, dbIDTracker.EventLogLastProcessed)
			tr.Duration = tr.BillrunEndTime.Sub(tr.BillingStartTime).String()
			// TODO
			tr.Done = true
		})
	}

	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventLogCurrent
	return nil
}

// BillingStarted checks if billing is started
func (c *Controller) BillingStarted(owner string, last uint64) bool {
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result = 'Starting Bill Run for owner ''%v'''", last, owner)
	var v uint32
	c.getLastVal(q, []interface{}{&v})

	if v != 0 {
		return true
	}
	return false
}

// BillingStartTime returns the billing start time
func (c *Controller) BillingStartTime(owner string, last uint64) time.Time {
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result = 'Starting Bill Run for owner ''%v'''", last, owner)
	var t time.Time
	c.getLastVal(q, []interface{}{&t})
	return t
}

// BillingFinished checks if billing is finished
func (c *Controller) BillingFinished(owner string, last uint64) bool {
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Billing for owner ''%v''%%'", last, owner)
	var v uint32
	c.getLastVal(q, []interface{}{&v})

	if v != 0 {
		return true
	}
	return false
}

// BillingEndTime returns the billing end time
func (c *Controller) BillingEndTime(owner string, last uint64) time.Time {
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Billing for owner ''%v''%%'", last, owner)
	var t time.Time
	c.getLastVal(q, []interface{}{&t})
	return t
}

// InvoiceRenderingStarted checks if invoice rendering is started
func (c *Controller) InvoiceRenderingStarted(owner string, last uint64) bool {
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Running Render Invoice for owner ''%v'''", last, owner)
	var v uint32
	c.getLastVal(q, []interface{}{&v})

	if v != 0 {
		return true
	}
	return false
}

// InvoiceRenderingStartTime returns invoice rendering start time
func (c *Controller) InvoiceRenderingStartTime(owner string, last uint64) time.Time {
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result = 'Running Render Invoice for owner ''%v'''", last, owner)
	var t time.Time
	c.getLastVal(q, []interface{}{&t})
	return t
}

// InvoiceRenderingFinished checks if invoice rendering is finished
func (c *Controller) InvoiceRenderingFinished(owner string, last uint64) bool {
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Render Invoice for owner ''%v'''", last, owner)
	var v uint32
	c.getLastVal(q, []interface{}{&v})

	if v != 0 {
		return true
	}
	return false
}

// InvoiceRenderingEndTime returns the invoice rendering end time
func (c *Controller) InvoiceRenderingEndTime(owner string, last uint64) time.Time {
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result = 'Finished Render Invoice for owner ''%v'''", last, owner)
	var t time.Time
	c.getLastVal(q, []interface{}{&t})
	return t
}

// BillrunFinished checks if billrun is completed
func (c *Controller) BillrunFinished(owner string, last uint64) bool {
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Bill Run for owner ''%v''%%'", last, owner)
	var v uint32
	c.getLastVal(q, []interface{}{&v})

	if v != 0 {
		return true
	}
	return false
}

// BillrunEndTime returns the bill run end time
func (c *Controller) BillrunEndTime(owner string, last uint64) time.Time {
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Bill Run for owner ''%v'''", last, owner)
	var t time.Time
	c.getLastVal(q, []interface{}{&t})
	return t
}
