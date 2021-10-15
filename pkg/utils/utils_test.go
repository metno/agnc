package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	firstDayOfYear := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	parsedDate := parseDate("NAD12301200123018001", firstDayOfYear)
	wantDate := time.Date(2020, 12, 30, 0, 0, 0, 0, time.UTC)
	if parsedDate != wantDate {
		t.Errorf("got: %v, expected: %v", parsedDate, wantDate)
	}

	parsedDate = parseDate("NAD01021200010218001", firstDayOfYear)
	wantDate = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	if parsedDate != wantDate {
		t.Errorf("got: %v, expected: %v", parsedDate, wantDate)
	}

	otherDate := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	parsedDate = parseDate("NAD02291200022918001", otherDate)
	wantDate = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	if parsedDate != wantDate {
		t.Errorf("got: %v, expected: %v", parsedDate, wantDate)
	}

	otherDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	parsedDate = parseDate("NAD02291200022918001", otherDate)
	wantDate = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	if parsedDate != wantDate {
		t.Errorf("got: %v, expected: %v", parsedDate, wantDate)
	}
}

func TestIsTooOldWrongFormat(t *testing.T) {
	testDate := time.Date(2021, 1, 17, 0, 0, 0, 0, time.UTC)
	_, err := IsTooOld("", testDate, 5)
	if err == nil {
		t.Errorf("This should have failed as it's not possible to extract the date from an empty string")
	}

}
func TestIsTooOld(t *testing.T) {
	testDate := time.Date(2021, 1, 17, 0, 0, 0, 0, time.UTC)
	iTO, _ := IsTooOld("NAD01161200011618001", testDate, 2)
	if iTO == true {
		t.Errorf("Instance should not be considered too old as it is not older than 2 days")
	}
}

func TestIsTooOldNot(t *testing.T) {
	testDate := time.Date(2021, 1, 20, 0, 0, 0, 0, time.UTC)
	iTO, _ := IsTooOld("NAD01161200011618001", testDate, 2)
	if iTO == false {
		t.Errorf("Instance should be considered too old as it is older than 2 days")
	}
}
