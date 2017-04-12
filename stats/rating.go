package stats

// GetRates collects the rates from the eventlog table
// returns the rates collected up to now, the number of files processed and the
// next id you should use for the next query
func (c *Controller) GetRates(filename string, lastID uint64) (updatedID uint64, filesProcessed int, result []float64) {
	return 0, 0, nil
	/*
		updatedID = lastID
		filesProcessed = 0

		var (
			filesCompletedRxp = regexp.MustCompile("Done Processing File" + ".*" + filename)
			InvalidRatesRxp   = regexp.MustCompile("UDRs in 0.0 seconds")
			RateRxp           = regexp.MustCompile("([0-9]+)*.([0-9]+)* UDRs/second|([0-9]+)* UDRs/second")
			RateValRxp        = regexp.MustCompile("([0-9]+)*.([0-9]+)*")
		)

		q := fmt.Sprintf("select id, result from "+
			"eventlog where id > %v and "+
			"(module = 'UDR Rating' or module = 'UDRRatingEngine') order by id", lastEventID)

		rows, err := c.db.Query(q)
		if err != nil {
			log.Println("WARNING: Stats controller generates an error while getting UDR rates: ", err)
			return updatedID, 0, nil
		}

		var id uint64
		var row string
		defer rows.Close()
		for rows.Next() {
			rowErr := rows.Scan(&id, &row)
			if rowErr != nil {
				log.Println("WARNING: Stats controller generates an error while scanning a row: ", err)
				return updatedID, 0, nil
			}

			// probably a overkill
			if updatedID < id {
				updatedID = id
			}

			if InvalidRatesRxp.MatchString(row) {
				continue
			}

			if filesCompletedRxp.MatchString(row) {
				filesProcessed++
			}

			if fs := RateRxp.FindString(row); fs != "" {
				fsv := RateValRxp.FindString(fs)

				r, err2 := strconv.ParseFloat(fsv, 64)
				if err2 == nil {
					result = append(result, r)
				}
			}
		}

		err = rows.Err()
		if err != nil {
			log.Println("WARNING: Stats controller generates an error: ", err)
			return updatedID, 0, nil
		}

		return updatedID, filesProcessed, result
	*/
}
