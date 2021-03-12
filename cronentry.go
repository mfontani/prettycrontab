package main

import (
	"fmt"
	"strings"
	"time"
)

// CronEntry represents one cron entry
type CronEntry struct {
	spec    string
	chunks  []string
	nextRun time.Time
	label   string
}

// CronEntries is the list of CronEntry iterated upon by the script
var CronEntries []CronEntry = make([]CronEntry, 0)

// SchedulingPart returns the scheduling part of the CronEntry
func (c CronEntry) SchedulingPart() string {
	return strings.Join(c.chunks[0:5], " ")
}

// CommandPart returns the command part of the CronEntry
func (c CronEntry) CommandPart() string {
	if len(c.label) > 0 {
		return c.label
	}
	if len(c.chunks) < 5 {
		return ""
	}
	cmdPart := strings.Join(c.chunks[5:], " ")
	if showRedirectDetails {
		return cmdPart
	}
	cmdWithoutRedirects := cmdPart
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, "2> /dev/null", "", -1)
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, "2>/dev/null", "", -1)
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, "1> /dev/null", "", -1)
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, "1>/dev/null", "", -1)
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, "> /dev/null", "", -1)
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, ">/dev/null", "", -1)
	cmdWithoutRedirects = strings.Replace(cmdWithoutRedirects, "2>&1", "", -1)
	return cmdWithoutRedirects
}

// Show the information about a CronEntry, depending on flags
func (c CronEntry) Show() string {
	chunks := make([]string, 0)

	if showDeltaSeconds {
		chunks = append(chunks, fmt.Sprintf("%d", c.nextRun.Unix()-now.Unix()))
	}
	if showDeltaCoarse {
		chunks = append(chunks, secondsToCoarseHMS(c.nextRun.Unix()-now.Unix()))
	} else if showDeltaHMS {
		chunks = append(chunks, secondsToHMS(c.nextRun.Unix()-now.Unix()))
	}
	if showTimestamp {
		chunks = append(chunks, c.nextRun.String())
	}
	if showCronSpec {
		chunks = append(chunks, c.SchedulingPart())
	}
	chunks = append(chunks, c.CommandPart())
	return strings.Join(chunks, "\t")
}
