package csvstore

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sonar/pkg/store"
	"strconv"
	"strings"
	"time"
)

func parseSonarRecord(fields []string) (store.SonarRecord, error) {
	// TODO: use a dictionary mapping desired keys to index in the fields to extend parsing to more CSV formats.
	record := store.SonarRecord{}

	if len(fields) != 4 {
		return record, errors.New("expected four field")
	}

	stamp, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return record, fmt.Errorf("failed to parse timestamp '%s' as int64: %s", fields[0], err)
	}
	record.Timestamp = time.Unix(stamp, 0)

	nums := [3]float32{}

	for index := range nums {
		value := strings.Trim(fields[index+1], " ")
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return record, fmt.Errorf("failed to parse number '%s': %s", value, err)
		}
		nums[index] = float32(num)
	}

	record.Longitude = nums[0]
	record.Latitude = nums[1]
	record.Depth = nums[2]

	return record, nil
}

// CreateMemoryStore parses a CSV content and creates an in-memory sonar data store.
//
// Expectation is that every CSV record is <timestamp>,<latitude>,<longitude>,<depth>
func CreateMemoryStore(reader io.Reader) (store.Memory, error) {
	// TODO use options pattern when we come across more nuances and settings.

	parser := csv.NewReader(reader)
	records := store.SonarRecordSet{}
	index := 0
	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}

		if index == 0 { // ignore header row
			index++
			continue
		}

		if err != nil {
			return store.Memory{}, err
		}

		// FIXME: use a separate parser here so if we get different CSV structure we can handle it.

		sonarRecord, err := parseSonarRecord(record)
		if err != nil {
			return store.Memory{}, fmt.Errorf("failed to parse record %v to sonar record: %s", record, err)
		}

		records = append(records, sonarRecord)
	}

	return store.Memory{Records: records}, nil
}

// CreateMemoryStoreFromFile creates an in-memory store from a CSV file.
func CreateMemoryStoreFromFile(path string) (store.Memory, error) {
	file, err := os.Open(path)
	if err != nil {
		return store.Memory{}, fmt.Errorf("failed to open CSV file: %s", err)
	}
	defer file.Close()

	return CreateMemoryStore(file)
}
