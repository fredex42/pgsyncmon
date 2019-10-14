package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

/**
numeric representation of a postgres WAL location
*/
type PostgresXlogLocation struct {
	Upper int32
	Lower int32
}

/**
calculate the difference between two WAL locations
*/
func (loc *PostgresXlogLocation) Difference(other PostgresXlogLocation) PostgresXlogLocation {
	return PostgresXlogLocation{
		Upper: loc.Upper - other.Upper,
		Lower: loc.Lower - other.Lower,
	}
}

/**
Utility class to make extracting WAL location data more efficient.
*/
type PostgresXlogLocationBuilder struct {
	xtractorRegex *regexp.Regexp
}

/**
Constructor for PostgresXlogLocationBuilder
*/
func NewPostgresXlogLocationBuilder() *PostgresXlogLocationBuilder {
	rtn := PostgresXlogLocationBuilder{
		xtractorRegex: regexp.MustCompile("^([0-9A-Fa-f]{4})/([0-9A-Fa-f]{8})$"),
	}
	return &rtn
}

/**
internal function to make parsing simpler
*/
func parseIntWrapper(str string) int32 {
	toConvert := fmt.Sprintf("0x%s", str)
	rtn, _ := strconv.ParseInt(toConvert, 0, 32)

	return int32(rtn)
}

/**
Call this method to get a PostgresXlogLocation from a string like
*/
func (b *PostgresXlogLocationBuilder) BuildFromString(str string) (*PostgresXlogLocation, error) {
	matches := b.xtractorRegex.FindStringSubmatch(str)
	if matches == nil {
		return nil, errors.New(fmt.Sprintf("%s did not match expected format for a postgres location", str))
	}

	rtn := PostgresXlogLocation{
		Upper: parseIntWrapper(matches[1]),
		Lower: parseIntWrapper(matches[2]),
	}
	return &rtn, nil
}

type RecoveryStatus struct {
	LastXlogReceive PostgresXlogLocation
	LastXlogReplay  PostgresXlogLocation
	IsInRecovery    bool
}

/**
calculate what our lag is, i.e. the difference between last chunk received and last chunk replayed
*/
func (r *RecoveryStatus) Lag() PostgresXlogLocation {
	return r.LastXlogReceive.Difference(r.LastXlogReplay)
}
