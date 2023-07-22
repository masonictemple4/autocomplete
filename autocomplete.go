// Autocomplete provides an in-memory autocompletion engine.
// It is designed to be used in conjunction with a data store to provide
// autocompletion suggestions for a given prefix. This is useful for things
// like search bars, command line interfaces, and more.
package autocomplete

import "time"

const SERVICE_NAME = "autocomplete"

type autocompleter interface {
	Insert(word string)
	Autocomplete(prefix string) []string
	Contains(word string) bool
	ListContents() []string
}

// Autocomplete service is the main object you will be interacting with.
// It is responsible for managing the autocompleter, data sources, and snapshots.
// It also provides a direct interface to interact with the autocompleter. That
// makes it easier than having to access the store to interface with the functionality.
type AutocompleteService struct {
	Config       Config
	SnapshotDest DataSource
	DataSources  []DataSource

	store autocompleter

	Errors      []error
	LastUpdated int64
	// TODO: Log
}

type Config struct {
	ServiceName      string
	MaxResults       int
	SnapshotsEnabled bool
	SnapshotInterval int

	AutomaticUpdates       bool
	LoadDataSourcesOnStart bool
	LowMemoryMode          bool
}

// New creates a new AutocompleteService instance and performs all of the setup.
// This makes a call to LoadDataSources(). If you wish to skip this,
// set the LoadDataSourcesOnStart option to false.
//
// You can also pass in a slice of keywords when calling this function to initialize
// your service store with.
func New(opts Config, keywords []string) (*AutocompleteService, error) {
	var store autocompleter
	if opts.LowMemoryMode {
		store = newTernarySearchTree("")
	} else {
		store = newTrie()
	}

	service := &AutocompleteService{
		Config: opts,
		store:  store,
		Errors: make([]error, 0),
	}

	for _, keyword := range keywords {
		service.store.Insert(keyword)
	}

	if opts.LoadDataSourcesOnStart {
		err := service.LoadDataSources()
		if err != nil {
			return nil, err
		}
	} else {
		// LoadDataSource will set the LastUpdated timestamp so we just
		// need to make sure if we don't call it we update it here.
		service.LastUpdated = time.Now().Unix()
	}

	return service, nil
}

func (a *AutocompleteService) LoadDataSources() error {
	for _, source := range a.DataSources {
		err := source.Provider.ReadData(source.Filepath, a.store, source.Formatter)
		if err != nil {
			a.Errors = append(a.Errors, err)
			return err
		}
	}
	a.LastUpdated = time.Now().Unix()
	return nil
}

func (a *AutocompleteService) CreateSnapshot() error {
	err := a.SnapshotDest.Provider.DumpData(a.SnapshotDest.Filepath, a.store, a.SnapshotDest.Formatter)
	if err != nil {
		a.Errors = append(a.Errors, err)
	}
	return err
}

func (a *AutocompleteService) RestoreFromSnapshot() error {
	err := a.SnapshotDest.Provider.ReadData(a.SnapshotDest.Filepath, a.store, a.SnapshotDest.Formatter)
	if err != nil {
		a.Errors = append(a.Errors, err)
		return err
	}
	a.LastUpdated = time.Now().Unix()
	return err
}

func (a *AutocompleteService) LoadDataSource(src DataSource) error {
	err := src.Provider.ReadData(src.Filepath, a.store, src.Formatter)
	if err != nil {
		a.Errors = append(a.Errors, err)
		return err
	}
	a.LastUpdated = time.Now().Unix()
	return nil
}

func (a *AutocompleteService) ExportToDataSource(dest DataSource) error {
	err := dest.Provider.DumpData(dest.Filepath, a.store, dest.Formatter)
	if err != nil {
		a.Errors = append(a.Errors, err)
		return err
	}
	return nil
}

// I am providing different names to these functions to avoid
// implementing the internal interface autocompleter on itself.
// This also provides quick access instead of having to go through
// the store. And gives us room to add more functionality later.
func (a *AutocompleteService) Complete(prefix string) []string {
	return a.store.Autocomplete(prefix)
}

func (a *AutocompleteService) Exists(word string) bool {
	return a.store.Contains(word)
}

func (a *AutocompleteService) Add(word string) {
	a.store.Insert(word)
}

func (a *AutocompleteService) GetContents(word string) []string {
	return a.store.ListContents()
}
