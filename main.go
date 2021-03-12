package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var showDeltaSeconds bool = false
var showDeltaHMS bool = false
var showDeltaCoarse bool = false
var showTimestamp bool = true
var showCronSpec bool = false
var showRedirectDetails bool = false

var now time.Time = time.Now()

func init() {
	flag.BoolVar(&showDeltaSeconds, "deltaseconds", showDeltaSeconds, `Show delta seconds until next run, before delta hms and full timestamp`)
	flag.BoolVar(&showDeltaHMS, "deltahms", showDeltaHMS, `Show delta h/m/s until next run, before full timestamp`)
	flag.BoolVar(&showDeltaCoarse, "deltacoarse", showDeltaCoarse, `Show "coarse" delta h/m/s until next run, before full timestamp`)
	flag.BoolVar(&showTimestamp, "timestamp", showTimestamp, `Show full timestamp, before cron spec`)
	flag.BoolVar(&showCronSpec, "spec", showCronSpec, `Show full/original cron spec (i.e. the "* * * * *" bits + command etc.).`)
	flag.BoolVar(&showRedirectDetails, "redir", showRedirectDetails, `(requires -spec=0) Show redirect details (">/dev/null", "2>&1") in command`)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\tcrontab -l | prettycrontab [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "Shows the next run time of interesting cron entries, from most recent to least.\n")
		fmt.Fprintf(os.Stderr, "You can tag a block of entries as `## UNINTERESTING` to have the block skipped.\n")
		fmt.Fprintf(os.Stderr, "You can tag an entry as `## LABEL foo` to have the label changed.\n")
		fmt.Fprintf(os.Stderr, "Otherwise, all commands are shown in the list.\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "\t* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # shown as-is\n")
		fmt.Fprintf(os.Stderr, "\t## UNINTERESTING\n\t* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # not shown\n")
		fmt.Fprintf(os.Stderr, "\t* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # not shown due to previous UNINTERESTING\n\n")
		fmt.Fprintf(os.Stderr, "\t* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # shown as-is as the double new line reset the block\n")
		fmt.Fprintf(os.Stderr, "\t## LABEL foo --bar\n\t* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # shown with custom label\n")
		fmt.Fprintf(os.Stderr, "\t* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # shown as-is as the custom label is only valid for one entry\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if showDeltaHMS && showDeltaCoarse {
		panic(fmt.Errorf("cannot -deltahms and -deltacoarse at the same time"))
	}
}

// Once finding "## UNINTERESTING", skip the next entries until an *empty line* is found.
var nextIsUninteresting bool = false

// Once finding "## LABEL ...", "label" the next *cron* entry as "label".
var foundLabel bool = false
var wantedLabel string = ""

// If the first chunk starts with a # then it's a comment!
var rxAllComments *regexp.Regexp = regexp.MustCompile(`\A\s*[#]`)

// If the first chunk looks like a variable setting
var rxVariableSet *regexp.Regexp = regexp.MustCompile(`\A\s*\w+\s*=`)

func main() {
	p := cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		spec := scanner.Text()
		chunks := strings.Fields(spec)
		// Skip empty lines
		if len(chunks) == 0 {
			nextIsUninteresting = false
			foundLabel = false
			continue
		}
		if chunks[0] == "##" && len(chunks) >= 2 && chunks[1] == "UNINTERESTING" {
			nextIsUninteresting = true
			continue
		}
		if chunks[0] == "##" && len(chunks) >= 2 && chunks[1] == "LABEL" {
			nextIsUninteresting = false
			foundLabel = true
			wantedLabel = strings.Join(chunks[2:], " ")
			continue
		}
		// Skip comment lines
		if rxAllComments.MatchString(chunks[0]) {
			continue
		}
		// Skip lines setting variables
		if rxVariableSet.MatchString(chunks[0]) {
			continue
		}
		if len(chunks) < 5 {
			panic(fmt.Errorf("Invalid cron entry too few fields: %s", spec))
		}
		cronPart := strings.Join(chunks[0:5], " ")
		s, err := p.Parse(cronPart)
		if err != nil {
			panic(fmt.Errorf("Could not parse cron part of %s: %s: %s", spec, cronPart, err))
		}
		if nextIsUninteresting {
			foundLabel = false
			continue
		}
		nextRun := s.Next(now)
		ce := CronEntry{
			spec:    spec,
			chunks:  chunks,
			nextRun: nextRun,
			label:   wantedLabel,
		}
		CronEntries = append(CronEntries, ce)
		wantedLabel = ""
		foundLabel = false
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	sort.Slice(CronEntries, func(i, j int) bool {
		return CronEntries[i].nextRun.Unix() < CronEntries[j].nextRun.Unix()
	})
	for _, ce := range CronEntries {
		fmt.Println(ce.Show())
	}
}
