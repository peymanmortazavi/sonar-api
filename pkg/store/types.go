package store

import (
	"errors"
	"time"
)

// Location demonstrates a geo location.
type Location struct {
	Longitude float32 `json:"longitude,omitempty"`
	Latitude  float32 `json:"latitude,omitempty"`
}

// SonarRecord holds a single sonar data measurement record.
type SonarRecord struct {
	Location  `json:",inline"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Depth     float32   `json:"depth,omitempty"`
}

// SonarRecordSet is a collection of sonar records.
type SonarRecordSet []SonarRecord

// Records simply returns all the records, a way to implement a Cursor.
func (s SonarRecordSet) Records() []SonarRecord {
	return s
}

// HasNext always returns false since a solar record collection is always final.
func (s SonarRecordSet) HasNext() bool {
	return false
}

// Next is not implemented for a sonar record set since it is a fixed storage.
func (s SonarRecordSet) Next() error {
	return errors.New("SonarRecordSet is a complete slice so Next is not implemented for this")
}

// Region describes a geo region.
type Region struct {
	NorthWest Location
	SouthEast Location
}

// Cursor describes a sonar data cursor.
type Cursor interface {
	Records() []SonarRecord
	HasNext() bool
	Next() error
}

// Reader describes the ability of reading all sonar records.
type Reader interface {
	Fetch(...FetchOption) (Cursor, error)
}

// Scanner describes the ability of searching/scanning all sonar records.
type Scanner interface {
	Scan(...ScanOption) (Cursor, error)
}

// ReadScanner describes a type that is both capable of returning records and querying them.
type ReadScanner interface {
	Reader
	Scanner
}
