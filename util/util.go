package util

import (
	"fmt"
	"time"
)

func PrintableDate(date time.Time) string {
	y, m, d := date.Date()
	return fmt.Sprintf("%d-%d-%d", y, m, d)
}

func PrintableTimeInCalling(daysInCalling int) string {
	years := daysInCalling / 365
	months := (daysInCalling - (365*years))/30
	days := daysInCalling - ((365 * years) + (30 * months))

	if months == 0 && days > 30 {
		months++
	}

	yearsLabel := "years"
	if years < 2 {
		yearsLabel = "year"
	}

	monthsLabel := "months"
	if months < 2 {
		monthsLabel = "month"
	}

	if years > 0 {
		if months > 0 {
			return fmt.Sprintf("%d %s, %d %s", years, yearsLabel, months, monthsLabel)
		}
		return fmt.Sprintf("%d %s", years, yearsLabel)
	}

	if months > 0 {
		return fmt.Sprintf("%d %s", months, monthsLabel)
	}

	return "A few days"
}
