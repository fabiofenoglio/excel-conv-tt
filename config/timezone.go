package config

import (
	"time"
	_ "time/tzdata"
)

const (
	fixedTimeZoneName = "Europe/Rome"
)

var (
	timeZone *time.Location
)

func init() {
	var err error
	timeZone, err = time.LoadLocation(fixedTimeZoneName)
	if err != nil {
		panic(err)
	}
}

func TimeZone() *time.Location {
	return timeZone
}
