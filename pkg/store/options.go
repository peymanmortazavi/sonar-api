package store

import "time"

// FetchOptions holds configurations of a fetch request.
type FetchOptions struct {
	BatchSize *int
}

// FetchOption describes a fetch option.
type FetchOption interface {
	Apply(*FetchOptions)
}

// Apply applies all given options to the current options.
func (f *FetchOptions) Apply(opts ...FetchOption) {
	for _, opt := range opts {
		opt.Apply(f)
	}
}

type simpleFetchOption func(*FetchOptions)

func (s simpleFetchOption) Apply(options *FetchOptions) {
	s(options)
}

// WithBatchSize sets the batch size of the returning cursor.
func WithBatchSize(size int) FetchOption {
	return simpleFetchOption(func(opts *FetchOptions) {
		opts.BatchSize = &size
	})
}

// ScanOptions holds all options and parameters related to a search.
type ScanOptions struct {
	FetchOptions
	Depth     *float32
	StartTime *time.Time
	EndTime   *time.Time
	Region    *Region
}

// Apply applies all given options to the current options.
func (s *ScanOptions) Apply(opts ...ScanOption) {
	for _, opt := range opts {
		opt.Apply(s)
	}
}

// ScanOption describes a scan option.
type ScanOption interface {
	Apply(*ScanOptions)
}

type simpleScanOption func(*ScanOptions)

func (s simpleScanOption) Apply(options *ScanOptions) {
	s(options)
}

// TODO: add an operation variable to allow for more granular search terms like greater than, equal to, less than, ...

// WithDepth only scans for depth greater than the specified depth.
func WithDepth(depth float32) ScanOption {
	return simpleScanOption(func(opts *ScanOptions) {
		opts.Depth = &depth
	})
}

// WithAfterTime only scans for records after the specified time.
func WithAfterTime(time time.Time) ScanOption {
	return simpleScanOption(func(opts *ScanOptions) {
		opts.StartTime = &time
	})
}

// WithBeforeTime only scans for records before the specified time.
func WithBeforeTime(time time.Time) ScanOption {
	return simpleScanOption(func(opts *ScanOptions) {
		opts.EndTime = &time
	})
}

// WithRegion only scans for records in the specified region.
func WithRegion(region Region) ScanOption {
	return simpleScanOption(func(opts *ScanOptions) {
		opts.Region = &region
	})
}

// WithFetchOptions applies fetching options for a scan request.
func WithFetchOptions(fetchOpts ...FetchOption) ScanOption {
	return simpleScanOption(func(opts *ScanOptions) {
		opts.FetchOptions.Apply(fetchOpts...)
	})
}
