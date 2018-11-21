package stats

import (
	"fmt"
	"log"
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

	// eventlog with "BillUDR" action completed count in eventlog
	wg.Add(1)
	go c.GetNumOfBillUDRActionsCompleted(&wg, dbIDTracker.EventlogStarted, &(tr.BillUDRActionCompleted))

	wg.Add(1)
	go c.GetNumOfUsageTranscationsGenerated(&wg, tp.OwnerName, dbIDTracker.EventlogStarted, &(tr.UsageTransactionsGenerated))

	wg.Add(1)
	go c.GetNumOfMRCTransactionsGenerated(&wg, tp.OwnerName, dbIDTracker.EventlogStarted, &(tr.MRCTransactionsGenerated))

	wg.Add(1)
	go c.GetNumOfInvoicesClosed(&wg, tp.OwnerName, dbIDTracker.EventlogStarted, &(tr.InvoicesClosed))

	wg.Wait()

	if !tr.BillingEndTime.IsZero() {
		if tr.BillingStartTime.IsZero() {
			panic("ERR: Billing end time captured but not start time")
		}

		tr.BillingEndTimeOnce.Do(func() {
			var wg sync.WaitGroup
			wg.Add(1)
			go c.GetMRCBillingDuration(&wg, dbIDTracker.EventlogStarted, &(tr.MRCBillingDuration))

			wg.Add(1)
			go c.GetUsageBillingDuration(&wg, dbIDTracker.EventlogStarted, &(tr.UsageBillingDuration))

			wg.Wait()
			// HAN >>>
			log.Printf("DEBUG: MRC billing duration: %v", tr.MRCBillingDuration)
			log.Printf("DEBUG: usage billing duration: %v", tr.UsageBillingDuration)
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
			panic("ERR: Invoice end time captured but not start time")
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

// GetNumOfBillUDRActionsCompleted gets the number of BillUDR actions with a "Finished Usage Billing for ..." keyword
func (c *Controller) GetNumOfBillUDRActionsCompleted(wg *sync.WaitGroup, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}
	var tp uint64
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and action = 'BillUDR' and Result like 'Finished Usage Billing for User%%'", last)
	c.getLastVal(q, []interface{}{tp})
	*result = append(*result, tp)
}

// GetDurationForAction calculates the time difference between the first entry and the last entry with the specified action in eventlog
func (c *Controller) GetDurationForAction(last uint64, action string) string {
	var startTime, endTime time.Time
	qEndTime := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = '%v' and Module = 'Billing' order by id desc", last, action)
	c.getLastVal(qEndTime, []interface{}{&endTime})
	// HAN >>>
	log.Printf("DEBUG: qEndTime query: %v", qEndTime)
	log.Printf("qEndTime: %v", endTime)
	qStartTime := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = '%v' and Module = 'Billing' order by id", last, action)
	c.getLastVal(qStartTime, []interface{}{&startTime})
	// HAN >>>
	log.Printf("DEBUG: qStartTime query: %v", qStartTime)
	log.Printf("qStartTime: %v", startTime)
	log.Printf("returns: %v", endTime.Sub(startTime).String())
	return endTime.Sub(startTime).String()
}

// GetUsageBillingDuration calculates the duration between the first eventlog with action "BillUDR" to the last
func (c *Controller) GetUsageBillingDuration(wg *sync.WaitGroup, last uint64, duration *string) {
	if wg != nil {
		defer wg.Done()
	}

	*duration = c.GetDurationForAction(last, "BillUDR")
	// HAN >>>
	log.Printf("DEBUG: duration: %v", duration)
}

// GetMRCBillingDuration calculates the duration between the first eventlog with action "BillUserPackage" to the last
func (c *Controller) GetMRCBillingDuration(wg *sync.WaitGroup, last uint64, duration *string) {
	if wg != nil {
		defer wg.Done()
	}

	*duration = c.GetDurationForAction(last, "BillUserPackage")
}

// BillrunStarted needs to return true to call stats gatherin functions below
// this is to make sure the billrun ID from billrunhistory table is
// indeed the bill run we are interested in
func (c *Controller) BillrunStarted(owner string, last uint64) bool {
	var startTime time.Time
	c.GetBillingStartTime(nil, owner, last, &startTime)
	return !startTime.IsZero()
}

// GetNumOfInvoicesClosed gets the number of closed invoices related to the latest billrun
func (c *Controller) GetNumOfInvoicesClosed(wg *sync.WaitGroup, owner string, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}

	// We have to correlate to eventlog to make sure that the last entry from
	// bill run history table is really the current billrun
	if !c.BillrunStarted(owner, last) {
		*result = append(*result, 0)
		return
	}

	var totalTrans uint64
	q := `select count(distinct(vi.InvoiceID))
				from viewinvoice vi
				inner join(select top 1 id from billrunhistory order by id desc) as A on vi.billrunhistoryid = A.id
				inner join viewstatementdetails vsd on vsd.InvoiceID = vi.InvoiceID
				inner join userpackage up on up.id = vsd.UserPackageID
				inner join invoicestatustype ist on ist.id = vi.InvoiceStatusTypeID
				where ist.name = 'Closed'`

	c.getLastVal(q, []interface{}{&totalTrans})
	*result = append(*result, totalTrans)
}

// GetNumOfTransactionsGenerated a helper function for both MRC and usage transactions
func (c *Controller) GetNumOfTransactionsGenerated(owner string, last uint64, isUsageTransaction bool, result *[]uint64) {
	// We have to correlate to eventlog to make sure that the last entry from
	// bill run history table is really the current billrun
	if !c.BillrunStarted(owner, last) {
		*result = append(*result, 0)
		return
	}

	var totalTrans uint64

	var whereClause string
	if isUsageTransaction {
		whereClause = `where vsd.name like '%Telecom Usage%'`
	} else {
		whereClause = `where vsd.name not like '%Telecom Usage%'`
	}

	q := `select count(vsd.id) from viewinvoice vi
	      inner join (select top 1 id from billrunhistory order by id desc) as A on vi.billrunhistoryid = A.id
				inner join viewstatementdetails vsd on vsd.InvoiceID = vi.InvoiceID ` + whereClause
	c.getLastVal(q, []interface{}{&totalTrans})
	// HAN >>>
	log.Printf("DEBUG: Usage transaction query: %v", q)
	log.Printf("DEBUG: Usage transaction generated: %v", totalTrans)
	*result = append(*result, totalTrans)
}

// GetNumOfUsageTranscationsGenerated gets the number of usage transactions related to the latest billrun
func (c *Controller) GetNumOfUsageTranscationsGenerated(wg *sync.WaitGroup, owner string, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}

	c.GetNumOfTransactionsGenerated(owner, last, true, result)
}

// GetNumOfMRCTransactionsGenerated gets the number of MRC transactions related to the latest billrun
func (c *Controller) GetNumOfMRCTransactionsGenerated(wg *sync.WaitGroup, owner string, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}

	c.GetNumOfTransactionsGenerated(owner, last, false, result)
}
