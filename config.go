package autocomplete

// ServiceConfig contains all of the configurable options for initializing a
// new autocomplete service.
//
// You can use the NewServiceConfig() function to create a new instance of this
type ServiceConfig struct {
	ServiceName string
	// Leave 0 for unlimited.
	MaxResults       int
	SnapshotsEnabled bool
	SnapshotInterval int

	AutomaticUpdates       bool
	LoadDataSourcesOnStart bool
	LowMemoryMode          bool

	SnapshotDest *DataSource
	DataSources  []DataSource
}

/* Config Functions */

// A type to help with a new pattern for passing options to the New() function.
type ConfigFn func(*ServiceConfig)

func WithServiceName(name string) ConfigFn {
	return func(c *ServiceConfig) {
		c.ServiceName = name
	}
}

// WithMaxResults sets the maximum number of results to return
// Leave this as 0 for unlimited.
func WithMaxResults(max int) ConfigFn {
	return func(c *ServiceConfig) {
		c.MaxResults = max
	}
}

func WithSnapshotsEnabled(c *ServiceConfig) {
	c.SnapshotsEnabled = true
}

func WithAutomaticUpdates(c *ServiceConfig) {
	c.AutomaticUpdates = true
}

func WithLoadDataSourcesOnStart(c *ServiceConfig) {
	c.LoadDataSourcesOnStart = true
}

func WithLowMemoryMode(c *ServiceConfig) {
	c.LowMemoryMode = true
}

func WithSnapshotInterval(interval int) ConfigFn {
	return func(c *ServiceConfig) {
		c.SnapshotInterval = interval
	}
}

func WithSnapshotDest(dest DataSource) ConfigFn {
	return func(c *ServiceConfig) {
		c.SnapshotDest = &dest
	}
}

func WithDataSources(sources []DataSource) ConfigFn {
	return func(c *ServiceConfig) {
		c.DataSources = sources
	}
}

/* End Config Functions */

// NewServiceConfig creates a new ServiceConfig instance with
// the default values. Then performs any updates based on the
// config functions passed in.
//
// Example:
//
//	config := NewServiceConfig(
//	  WithServiceName("my-service"),
//	  WithMaxResults(10),
//	  WithSnapshotsEnabled(),
//	)
//
// Will create a new ServiceConfig with the default values,
// then update Service name, MaxResults, and enable snapshots.
func NewServiceConfig(opts ...ConfigFn) *ServiceConfig {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func defaultConfig() *ServiceConfig {
	d, err := NewLocalFileProvider("/var/tmp/autocomplete/snapshot.json")
	if err != nil {
		panic(err)
	}

	snapshotDest := NewDataSource(d, DefaultFormat{}, d.Filename, "")
	return &ServiceConfig{
		ServiceName:            SERVICE_NAME,
		MaxResults:             0,
		SnapshotsEnabled:       false,
		SnapshotInterval:       0,
		AutomaticUpdates:       false,
		LoadDataSourcesOnStart: false,
		LowMemoryMode:          false,

		SnapshotDest: snapshotDest,
	}

}
