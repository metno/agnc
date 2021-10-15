package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func parseDate(ecFileName string, onDate time.Time) time.Time {
	year := onDate.Year()
	monthInt, _ := strconv.Atoi(ecFileName[3:5])
	month := time.Month(monthInt)
	day, _ := strconv.Atoi(ecFileName[5:7])
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if t.After(onDate) {
		t = time.Date(year-1, month, day, 0, 0, 0, 0, time.UTC)
	}

	return t
}

// IsTooOld function checks whether the grib file is older than ageInDays based on its name
func IsTooOld(ecFileName string, onDate time.Time, ageInDays int) (bool, error) {
	if len(ecFileName) != 20 {
		return false, errors.New(fmt.Sprintf("Only strings of length 20 are valid as ecFileName. Name received: %s", ecFileName))
	}

	expiresdAt := parseDate(ecFileName, onDate).AddDate(0, 0, ageInDays)

	if expiresdAt.Before(onDate) {
		return true, nil
	}

	return false, nil
}
