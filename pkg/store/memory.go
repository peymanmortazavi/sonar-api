package store

import (
	"errors"
	"time"
)

// Memory implements a simple in-memory store for sonar data.
type Memory struct {
	Records SonarRecordSet
}

type memoryCursor struct {
	records   []SonarRecord
	index     int
	batchSize int
}

func (m *memoryCursor) Records() []SonarRecord {
	return m.records[m.index : m.batchSize+m.index]
}

func (m *memoryCursor) HasNext() bool {
	return (m.index + m.batchSize) > len(m.records)
}

func (m *memoryCursor) Next() error {
	if !m.HasNext() {
		return errors.New("the cursor is exhusted")
	}
	m.index = m.index + m.batchSize
	return nil
}

// Fetch returns all the records in the in-memory store all at once.
func (m Memory) Fetch(opts ...FetchOption) (Cursor, error) {
	options := &FetchOptions{}
	options.Apply(opts...)

	if options.BatchSize != nil {
		return &memoryCursor{
			records:   m.Records,
			index:     0,
			batchSize: *options.BatchSize,
		}, nil
	}
	return m.Records, nil
}

// query describes a behavior to select or skip a sonar record.
type query interface {
	ShouldKeep(*SonarRecord) bool
}

type queryFunc func(*SonarRecord) bool

func (q queryFunc) ShouldKeep(r *SonarRecord) bool {
	return q(r)
}

func belowDepthQuery(depth float32) queryFunc {
	return queryFunc(func(s *SonarRecord) bool {
		return s.Depth > depth
	})
}

func afterTimeQuery(t time.Time) queryFunc {
	return queryFunc(func(s *SonarRecord) bool {
		return s.Timestamp.After(t)
	})
}

func beforeTimeQuery(t time.Time) queryFunc {
	return queryFunc(func(s *SonarRecord) bool {
		return s.Timestamp.Before(t)
	})
}

func regionQuery(region Region) queryFunc {
	return queryFunc(func(s *SonarRecord) bool {
		return (s.Location.Latitude > region.NorthWest.Latitude &&
			s.Location.Latitude < region.SouthEast.Latitude &&
			s.Location.Longitude > region.NorthWest.Latitude &&
			s.Location.Longitude < region.SouthEast.Latitude)
	})
}

func find(s SonarRecordSet, querySet ...query) SonarRecordSet {
	records := SonarRecordSet{}

	// if no query is specified, it is easy, return the entire set!
	if len(querySet) == 0 {
		return s
	}

	// apply all the query matchers!
mainLoop:
	for _, record := range s {
		for _, q := range querySet {
			if !q.ShouldKeep(&record) {
				continue mainLoop
			}
		}
		records = append(records, record)
	}

	return records
}

// Scan returns all matching records according to the search criteria.
func (m Memory) Scan(opts ...ScanOption) (Cursor, error) {
	// FIXME: actually support batching instead of running the query across the entire slice.
	options := &ScanOptions{}
	options.Apply(opts...)

	querySet := []query{}
	if depth := options.Depth; depth != nil {
		querySet = append(querySet, belowDepthQuery(*depth))
	}

	if startTime := options.StartTime; startTime != nil {
		querySet = append(querySet, afterTimeQuery(*startTime))
	}

	if endTime := options.EndTime; endTime != nil {
		querySet = append(querySet, beforeTimeQuery(*endTime))
	}

	if region := options.Region; region != nil {
		querySet = append(querySet, regionQuery(*region))
	}

	return find(m.Records, querySet...), nil
}
