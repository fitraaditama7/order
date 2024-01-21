package utils

import (
	"order-backend/constant"
	"time"
)

const (
	layoutDate = "2006-01-02"
)

func FormatStartDate(dateStr string) (*time.Time, error) {
	return formatDateWithTime(dateStr, false)
}

func FormatEndDate(dateStr string) (*time.Time, error) {
	return formatDateWithTime(dateStr, true)
}

func formatDateWithTime(dateStr string, isEndDate bool) (*time.Time, error) {
	location, err := time.LoadLocation(constant.MELBOURNE_LOCATION)
	if err != nil {
		return nil, err
	}

	t, err := time.ParseInLocation(layoutDate, dateStr, location)
	if err != nil {
		return nil, err
	}

	if isEndDate {
		t = t.AddDate(0, 0, 1)
	}

	// Format the date
	return &t, nil
}
