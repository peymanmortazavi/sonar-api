package restapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sonar/pkg/store"
	"strconv"
	"strings"
	"time"
)

// Handler is a simple HTTP handler for an sonar data store.
type Handler struct {
	store.ReadScanner
}

func createScanOptions(request *http.Request) ([]store.ScanOption, error) {
	opts := []store.ScanOption{}
	params := request.URL.Query()

	log.Printf("%v", params)

	floatFields := []string{"depth"}
	floats := map[string]float32{}

	for _, field := range floatFields {
		if fieldValue := params.Get(field); len(fieldValue) != 0 {
			value, err := strconv.ParseFloat(fieldValue, 32)
			if err != nil {
				return opts, fmt.Errorf("failed to parse %s '%s': %s", field, fieldValue, err)
			}
			floats[field] = float32(value)
		}
	}

	if depth, ok := floats["depth"]; ok {
		opts = append(opts, store.WithDepth(depth))
	}

	longFields := []string{"start", "end"}
	longs := map[string]int64{}

	for _, field := range longFields {
		if fieldValue := params.Get(field); len(fieldValue) != 0 {
			value, err := strconv.ParseInt(fieldValue, 10, 64)
			if err != nil {
				return opts, fmt.Errorf("failed to parse %s '%s': %s", field, fieldValue, err)
			}
			longs[field] = value
		}
	}

	if after, ok := longs["after"]; ok {
		opts = append(opts, store.WithAfterTime(time.Unix(after, 0)))
	}

	if before, ok := longs["before"]; ok {
		opts = append(opts, store.WithBeforeTime(time.Unix(before, 0)))
	}

	if regionValue := params.Get("region"); len(regionValue) != 0 {
		parts := strings.Split(regionValue, ",")
		if len(parts) != 4 {
			return opts, fmt.Errorf("expected 4 comma separated values but got '%s'", regionValue)
		}
		nums := [4]float32{}
		for index := range nums {
			value, err := strconv.ParseFloat(parts[index], 32)
			if err != nil {
				return opts, fmt.Errorf("failed to parse '%s': %s", parts[index], err)
			}
			nums[index] = float32(value)
		}
		opts = append(opts, store.WithRegion(store.Region{
			NorthWest: store.Location{Longitude: nums[0], Latitude: nums[1]},
			SouthEast: store.Location{Longitude: nums[2], Latitude: nums[3]},
		}))
	}

	return opts, nil
}

type recordResult struct {
	Records store.SonarRecordSet `json:"records"`
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// TODO: create a more generic error so it can be parsed out by downstream components.
	// TODO: support JSON and other formats by looking at the Content-Type header.
	options, err := createScanOptions(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	records, err := h.Scan(options...)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "failed to write query data: %s", err)
		return
	}

	if err := json.NewEncoder(writer).Encode(recordResult{records.Records()}); err != nil {
		log.Printf("failed to write result: %s", err)
	}
}
