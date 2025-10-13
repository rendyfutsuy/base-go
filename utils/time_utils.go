package utils

import (
	"fmt"
	"time"
)

// MonthlyIntervalCount calculates the number of months between two dates
func MonthlyIntervalCount(start, end time.Time) float64 {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	totalMonths := years*12 + months

	return float64(totalMonths)
}

// QuarterlyIntervalCount calculates the number of quarters between two dates
func QuarterlyIntervalCount(start, end time.Time) float64 {
	monthlyCount := MonthlyIntervalCount(start, end)
	return monthlyCount / 3
}

// HalfYearlyIntervalCount calculates the number of half-year intervals between two dates
func HalfYearlyIntervalCount(start, end time.Time) float64 {
	monthlyCount := MonthlyIntervalCount(start, end)
	return monthlyCount / 6
}

// check if the time t is between [min, max]
func TimeIsBetween(t, min, max time.Time) bool {
	if min.After(max) {
		min, max = max, min
	}
	return (t.Equal(min) || t.After(min)) && (t.Equal(max) || t.Before(max))
}

// UniversalTimeParser parse date with every possible format
func UniversalTimeParser(source string) (*time.Time, error) {

	date, err := time.Parse("2006-01-02T15:04:05Z", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2006-01-02T00:00:00", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2006-01-02 15:04:05", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2006-01-02", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2-Jan-06", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2-Jan-2006", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02/01/2006", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("01-02-06", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02-01-2006", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02.01.2006", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02 01 2006", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02-01-06", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02.01.06", source)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse("02 01 06", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.ANSIC, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.UnixDate, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RubyDate, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC822, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC822Z, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC850, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC1123, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC1123Z, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC3339, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC3339Nano, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2006-01-02T15:04:05.000Z", source)
	if err == nil {
		return &date, nil
	}

	return nil, fmt.Errorf("Requested String is not in Time format")
}

func GetDateDiff(a, b time.Time) (years, months, days int) {
	// ensure a is before b
	if a.After(b) {
		a, b = b, a
	}

	// extract date
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()

	// process
	years = yearB - yearA
	months = int(monthB - monthA)
	days = dayB - dayA

	if days < 0 {
		// days from previous month
		prevMonth := b.AddDate(0, -1, 0)
		_, _, daysInPrevMonth := prevMonth.Date()
		lastDayOfPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month(), daysInPrevMonth, 0, 0, 0, 0, b.Location()).Day()
		days += lastDayOfPrevMonth
		months--
	}

	if months < 0 {
		months += 12
		years--
	}

	return
}

func DateValidate(t time.Time) bool {
	var result bool

	if t.Year() == 1 && t.Month() == time.January && t.Day() == 1 {
		result = true
	}

	return result
}

func DefaultFormat() string {
	return "02 January 2006"
}
