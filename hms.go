package main

import (
	"fmt"
	"strings"
)

func secondsToHMS(seconds int64) string {
	s := seconds % 60
	m := int(seconds/60) % 60
	h := int(seconds/60/60) % 24
	d := int(seconds / 60 / 60 / 24)
	M := int(seconds / 60 / 60 / 24 / 30)
	if M != 0 {
		d = d - M*30
	}
	var outStrings []string
	if M != 0 {
		outStrings = append(outStrings, fmt.Sprintf("[%dM]", M))
	}
	if d != 0 {
		outStrings = append(outStrings, fmt.Sprintf("%dd", d))
	}
	if h != 0 {
		outStrings = append(outStrings, fmt.Sprintf("%dh", h))
	}
	if m != 0 {
		outStrings = append(outStrings, fmt.Sprintf("%dm", m))
	}
	if s != 0 {
		outStrings = append(outStrings, fmt.Sprintf("%ds", s))
	}
	// outStrings = append(outStrings, fmt.Sprintf("= %d seconds", seconds))
	return strings.Join(outStrings, " ")
}

func secondsToCoarseHMS(seconds int64) string {
	// s := seconds % 60
	m := int(seconds/60) % 60
	h := int(seconds/60/60) % 24
	d := int(seconds / 60 / 60 / 24)
	M := int(seconds / 60 / 60 / 24 / 30)
	if M != 0 {
		d = d - M*30
	}
	if M != 0 {
		if M > 6 {
			return secondsToHMS(seconds)
		}
		if M > 1 {
			return fmt.Sprintf("%dM", M)
		}
		return "a month"
	}
	if d != 0 {
		if d > 1 {
			return fmt.Sprintf("%dd", d)
		}
		return "a day"
	}
	if h != 0 {
		if h > 1 {
			return fmt.Sprintf("%dhr", h)
		}
		return "an hour"
	}
	if m > 10 {
		return fmt.Sprintf("%d mins", m)
	}
	return secondsToHMS(seconds)
}
