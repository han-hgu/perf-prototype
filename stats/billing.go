package stats

import (
	"fmt"
	"log"
	"strconv"
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
	go c.GetNumOfUsageTranscationsGenerated(&wg, tp.OwnerName, dbIDTracker.StatementDetailsStarted, &(tr.UsageTransactionsGenerated))

	wg.Add(1)
	go c.GetNumOfMRCTransactionsGenerated(&wg, tp.OwnerName, dbIDTracker.StatementDetailsStarted, &(tr.MRCTransactionsGenerated))

	wg.Add(1)
	go c.GetNumOfInvoicesClosed(&wg, tp.OwnerName, dbIDTracker.InvoiceStarted, &(tr.InvoicesClosed))

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
			// summarize duration for all actions with billing module
			ret, _ := c.GetBillingActions(dbIDTracker.EventlogStarted)
			tr.ActionDuration = ret

			tr.Duration = tr.BillrunEndTime.Sub(tr.BillingStartTime).String()
			tr.Done = true
		})
	}

	return nil
}

// GetBillingStartTime gets the billing start time from the latest event log entries
func (c *Controller) GetBillingStartTime(wg *sync.WaitGroup, owner string, last uint64, billingStartTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Module = 'Billing' and Result like 'Starting Bill Run%%for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{billingStartTime})
}

// GetBillingEndTime gets the billing end time from the latest event log entries
func (c *Controller) GetBillingEndTime(wg *sync.WaitGroup, owner string, last uint64, billingEndTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Module = 'Billing' and Result like 'Finished Billing for owner ''%v''%%' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{billingEndTime})
}

// GetInvoiceRenderStartTime gets the invoice render start time from the latest event log entries
func (c *Controller) GetInvoiceRenderStartTime(wg *sync.WaitGroup, owner string, last uint64, invoiceRenderStartTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Module = 'Billing' and Result = 'Running Render Invoice for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{invoiceRenderStartTime})
}

// GetInvoiceRenderEndTime gets the invoice render end time from the latest event log entries
func (c *Controller) GetInvoiceRenderEndTime(wg *sync.WaitGroup, owner string, last uint64, invoiceRenderEndTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Module = 'Billing' and Result = 'Finished Render Invoice for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{invoiceRenderEndTime})
}

// GetBillrunEndTime gets the bill run time from the latest event log entries
func (c *Controller) GetBillrunEndTime(wg *sync.WaitGroup, owner string, last uint64, billrunEndTime *time.Time) {
	if wg != nil {
		defer wg.Done()
	}
	q := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Module = 'Billing' and Result like 'Finished Bill Run%%for owner ''%v''' order by id desc", last, owner)
	c.getLastVal(q, []interface{}{billrunEndTime})
}

// GetNumOfBillUDRActionsCompleted gets the number of BillUDR actions with a "Finished Usage Billing for ..." keyword
func (c *Controller) GetNumOfBillUDRActionsCompleted(wg *sync.WaitGroup, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}
	var tp uint64
	q := fmt.Sprintf("select count(*) from eventlog where id > %v and action = 'BillUDR' and Result like 'Finished Usage Billing for User%%'", last)
	c.getLastVal(q, []interface{}{&tp})
	*result = append(*result, tp)
}

// GetDurationForAction calculates the time difference between the first entry and the last entry with the specified action in eventlog
func (c *Controller) GetDurationForAction(last uint64, action string) string {
	var startTime, endTime time.Time
	qEndTime := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = '%v' and Module = 'Billing' order by id desc", last, action)
	c.getLastVal(qEndTime, []interface{}{&endTime})
	qStartTime := fmt.Sprintf("select top 1 Date from eventlog where id > %v and Action = '%v' and Module = 'Billing' order by id", last, action)
	c.getLastVal(qStartTime, []interface{}{&startTime})
	return endTime.Sub(startTime).String()
}

// Get BillingActions groups the action from eventlog and report the number of entries, the duration between the first and the last and
// the rate if duration is not 0
func (c *Controller) GetBillingActions(last uint64) (map[string]map[string]interface{}, error) {
	ret := make(map[string]map[string]interface{})

	rows, err := c.db.Query("select action, datediff(s, min(date), max(date)), count(1) " +
		"from eventlog " +
		"where module = 'Billing' and " + fmt.Sprintf("id > %v ", last) +
		"group by action")

	var action string
	var duration, itemCount uint64
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&action, &duration, &itemCount)
		if err != nil {
			return make(map[string]map[string]interface{}), err
		}
		currMap := make(map[string]interface{})
		var dur time.Duration
		dur, _ = time.ParseDuration(strconv.FormatUint(duration, 10) + "s")
		currMap["item_count"] = itemCount
		currMap["duration"] = dur.String()
		if duration == 0 {
			currMap["rate"] = 0
		} else {
			currMap["rate"] = float64(itemCount) / float64(duration)
		}
		ret[action] = currMap
	}

	return ret, nil
}

// GetUsageBillingDuration calculates the duration between the first eventlog with action "BillUDR" to the last
func (c *Controller) GetUsageBillingDuration(wg *sync.WaitGroup, last uint64, duration *string) {
	if wg != nil {
		defer wg.Done()
	}

	*duration = c.GetDurationForAction(last, "BillUDR")
}

// GetMRCBillingDuration calculates the duration between the first eventlog with action "BillUserPackage" to the last
func (c *Controller) GetMRCBillingDuration(wg *sync.WaitGroup, last uint64, duration *string) {
	if wg != nil {
		defer wg.Done()
	}

	*duration = c.GetDurationForAction(last, "DoBillUserPackage")
}

// GetNumOfInvoicesClosed gets the number of invoices closed within the interval
func (c *Controller) GetNumOfInvoicesClosed(wg *sync.WaitGroup, owner string, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}

	var invoicesClosed uint64
	q := `select count(*) from invoice
	inner join InvoiceStatusType ist on ist.id = invoice.InvoiceStatusTypeID
	where ist.name = 'Closed' ` + fmt.Sprintf("and invoice.id > %v", last)

	c.getLastVal(q, []interface{}{&invoicesClosed})
	log.Printf("DEBUG: invoices closed query: %v", q)
	log.Printf("DEBUG: invoices closed in total: %v", invoicesClosed)
	*result = append(*result, invoicesClosed)
}

// GetNumOfTransactionsGenerated a helper function for both MRC and usage transactions
func (c *Controller) GetNumOfTransactionsGenerated(owner string, last uint64, isUsageTransaction bool, result *[]uint64) {
	var totalTrans uint64

	var whereClause string
	if isUsageTransaction {
		whereClause = `where detail like '%Telecom Usage%' `
	} else {
		whereClause = `where detail not like '%Telecom Usage%' `
	}

	whereClause += fmt.Sprintf("and id > %v", last)

	q := "select count(*) from statementdetails " + whereClause
	c.getLastVal(q, []interface{}{&totalTrans})
	*result = append(*result, totalTrans)
}

// GetNumOfUsageTranscationsGenerated gets the number of usage transactions up to now
func (c *Controller) GetNumOfUsageTranscationsGenerated(wg *sync.WaitGroup, owner string, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}

	c.GetNumOfTransactionsGenerated(owner, last, true, result)
}

// GetNumOfMRCTransactionsGenerated gets the number of MRC transactions up to now
func (c *Controller) GetNumOfMRCTransactionsGenerated(wg *sync.WaitGroup, owner string, last uint64, result *[]uint64) {
	if wg != nil {
		defer wg.Done()
	}

	c.GetNumOfTransactionsGenerated(owner, last, false, result)
}
