package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func StringToUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	return uint(i), err
}

func CountNumberOfFieldsInAStruct(obj interface{}) int {
	fields := reflect.VisibleFields(reflect.TypeOf(obj))
	count := 0
	for _, field := range fields {
		if !field.Anonymous {
			count += 1
		}
	}
	return count
}

func DurationToString(d time.Duration) string {
	d = d.Round(time.Second)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func StringToDuration(d string) *time.Duration {
	duration, err := time.ParseDuration(d)
	if err != nil {
		duration, _ := time.ParseDuration("0m")
		return &duration
	}
	return &duration
}

func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
