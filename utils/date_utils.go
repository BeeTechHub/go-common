package utils

import (
	"time"
	_ "time/tzdata"
)

const GmtPlus7 string = "Asia/Bangkok"
const DTFM_DATE_TIME_01 string = "200601021504"
const DTFM_DATE_TIME_02 string = "20060102150405"
const DTFM_DATE_TIME_03 string = "02/01/2006 15:04"
const DTFM_DATE_TIME_04 string = "02/01/2006 15:04:05"
const DTFM_DATE_TIME_05 string = "2006-01-02T15:04:05Z"
const DTFM_DATE_01 string = "20060102"
const DTFM_DATE_02 string = "02/01/2006"
const DTFM_DATE_MONTH_01 string = "200601"

// Time to string format "dd/MM/yyyy HH:mm" in GMT+7
func Time_ToStringDateTime03_GmtPlus7(input time.Time) (string, error) {
	location, err := time.LoadLocation(GmtPlus7)
	if err != nil {
		return "", err
	}

	// Convert the current time to the desired timezone
	nowInGmtPlus7 := input.In(location)

	// Format the time using the desired layout
	formattedTime := nowInGmtPlus7.Format(DTFM_DATE_TIME_03)

	return formattedTime, nil
}

// Time to string format "yyyyMMdd" in GMT+7
func Time_ToStringDate01_GmtPlus7(input time.Time) (string, error) {
	location, err := time.LoadLocation(GmtPlus7)
	if err != nil {
		return "", err
	}

	// Convert the current time to the desired timezone
	nowInGmtPlus7 := input.In(location)

	// Format the time using the desired layout
	formattedTime := nowInGmtPlus7.Format(DTFM_DATE_01)

	return formattedTime, nil
}

// Time to string format "yyyyMM" in GMT+7
func Time_ToStringDateMonth01_GmtPlus7(input time.Time) (string, error) {
	location, err := time.LoadLocation(GmtPlus7)
	if err != nil {
		return "", err
	}

	// Convert the current time to the desired timezone
	nowInGmtPlus7 := input.In(location)

	// Format the time using the desired layout
	formattedTime := nowInGmtPlus7.Format(DTFM_DATE_MONTH_01)

	return formattedTime, nil
}

// String format "2006-01-02T15:04:05Z" to time
func StringRFC3339Nano_ToTime(input string) (time.Time, error) {
	t, err := time.Parse(DTFM_DATE_TIME_05, input)
	return t, err
}

// String format "yyyyMMdd" in GMT+0 to time
func StringDate01_ToTime(input string) (time.Time, error) {
	t, err := time.Parse(DTFM_DATE_01, input)
	return t, err
}

// string format "yyyyMMdd" in GMT+7 to time
func StringDate01_GmtPlus7_ToTime(input string) (time.Time, error) {
	location, err := time.LoadLocation(GmtPlus7) // GMT+7
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(DTFM_DATE_01, input, location)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
