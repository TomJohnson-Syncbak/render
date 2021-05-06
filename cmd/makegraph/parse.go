package main

import "strconv"
import "time"

func parseTest() {
	parseTime("1")
	parseTime("7h25m42s")
	parseTime("1d")
	parseTime("1w")
	parseTime("h")
}

func parseTime(s string) (d time.Duration, err error) {
	i, err := strconv.Atoi(s)
	if err == nil {
		d = time.Duration(i) * time.Hour
		return
	}
	// Could not convert string to integer, try Duration now
	d, err = time.ParseDuration(s)
	return
}
