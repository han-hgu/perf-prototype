package stats

import (
	"fmt"
	"sync"
	"time"

	"github.com/perf-prototype/perftest"
)

// UpdateBillingResult to update the billing results
func (c *Controller) UpdateBillingResult(ti *perftest.TestInfo, dbIDTracker *perftest.DBIDTracker) error {
	dbIDTracker.EventLogCurrent = c.getLastEventLogID()
	tr := ti.Result.(*perftest.BillingResult)
	tp := ti.Params.(*perftest.BillingParams)

	var wg sync.WaitGroup

	if tr.BillingStartTime.IsZero() {
		wg.Add(1)
		go c.GetBillingStartTime(&wg, tp.OwnerName, dbIDTracker.EventLogLastProcessed, &(tr.BillingStartTime))
	}

	if tr.BillingEndTime.IsZero() {
		wg.Add(1)
		go c.GetBillingEndTime(&wg, tp.OwnerName, dbIDTracker.EventLogLastProcessed, &(tr.BillingEndTime))
	}

	if tr.InvoiceRenderStartTime.IsZero() {
		wg.Add(1)
		go c.GetInvoiceRenderStartTime(&wg, tp.OwnerName, dbIDTracker.EventLogLastProcessed, &(tr.InvoiceRenderStartTime))
	}

	if tr.InvoiceRenderEndTime.IsZero() {
		wg.Add(1)
		go c.GetInvoiceRenderEndTime(&wg, tp.OwnerName, dbIDTracker.EventLogLastProcessed, &(tr.InvoiceRenderEndTime))
	}

	if tr.BillrunEndTime.IsZero() {
		wg.Add(1)
		go c.GetBillrunEndTime(&wg, tp.OwnerName, dbIDTracker.EventLogLastProcessed, &(tr.BillrunEndTime))
	}

	wg.Wait()

	if !tr.BillingEndTime.IsZero() {
		if tr.BillingStartTime.IsZero() {
			panic("ERR: Billing end time captured but not start time")
		}

		tr.BillingEndTimeOnce.Do(func() {
			tr.BillingDuration = tr.BillingEndTime.Sub(tr.BillingStartTime).String()
		})
	}

	if !tr.InvoiceRenderEndTime.IsZero() {
		if tr.InvoiceRenderStartTime.IsZero() {
			panic("ERR: Invoice render end time captured but not start time")
		}

		tr.InvoiceRenderEndTimeOnce.Do(func() {
			tr.InvoiceRenderDuration = tr.InvoiceRenderEndTime.Sub(tr.InvoiceRenderStartTime).String()
		})
	}

	if !tr.BillrunEndTime.IsZero() {
		if tr.BillingStartTime.IsZero() {
			panic("ERR: Invo end time captured but not start time")
		}

		tr.BillrunEndOnce.Do(func() {
			tr.Duration = tr.BillrunEndTime.Sub(tr.BillingStartTime).String()
			tr.Done = true
		})
	}

	dbIDTracker.EventLogLastProcessed = dbIDTracker.EventLogCurrent
	return nil
}

// GetBillingStartTime gets the billing start time from the latest event log entries
func (c *Controller) GetBillingStartTime(wg *sync.WaitGroup, owner string, last uint64, billingStartTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Starting Bill Run%%for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{billingStartTime})
}

// GetBillingEndTime gets the billing end time from the latest event log entries
func (c *Controller) GetBillingEndTime(wg *sync.WaitGroup, owner string, last uint64, billingEndTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Billing for owner ''%v''%%' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{billingEndTime})
}

// GetInvoiceRenderStartTime gets the invoice render start time from the latest event log entries
func (c *Controller) GetInvoiceRenderStartTime(wg *sync.WaitGroup, owner string, last uint64, invoiceRenderStartTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result = 'Running Render Invoice for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{invoiceRenderStartTime})
}

// GetInvoiceRenderEndTime gets the invoice render end time from the latest event log entries
func (c *Controller) GetInvoiceRenderEndTime(wg *sync.WaitGroup, owner string, last uint64, invoiceRenderEndTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result = 'Finished Render Invoice for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{invoiceRenderEndTime})
}

// GetBillrunEndTime gets the bill run time from the latest event log entries
func (c *Controller) GetBillrunEndTime(wg *sync.WaitGroup, owner string, last uint64, billrunEndTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = 'CheckForBillRun' and Module = 'Billing' and Result like 'Finished Bill Run%%for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{billrunEndTime})
}
