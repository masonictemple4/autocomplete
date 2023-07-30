// Autocomplete provides an in-memory autocompletion engine.
// It is designed to be used in conjunction with a data store to provide
// autocompletion suggestions for a given prefix. This is useful for things
// like search bars, command line interfaces, and more.
package autocomplete

import (
	"fmt"
	"runtime"
	"time"
)

const SERVICE_NAME = "autocomplete"

type autocompleter interface {
	Insert(word string)
	Autocomplete(prefix string) []string
	Contains(word string) bool
	ListContents() []string
	Clear()
}

// Autocomplete service is the main object you will be interacting with.
// It is responsible for managing the autocompleter, data sources, and snapshots.
// It also provides a direct interface to interact with the autocompleter. That
// makes it easier than having to access the store to interface with the functionality.
type AutocompleteService struct {
	Config ServiceConfig

	store autocompleter

	Errors      []error
	LastUpdated int64
	isClosed    bool
	// TODO: Log
}

// New creates a new AutocompleteService instance and performs all of the setup.
// This makes a call to LoadDataSources(). If you wish to skip this,
// set the LoadDataSourcesOnStart option to false.
//
// You can also pass in a slice of keywords when calling this function to initialize
// your service store with.
func New(opts ServiceConfig, keywords []string) (*AutocompleteService, error) {
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

// Close will check for the SnapshotDest, and DataSources and close
// the providers associated with each. This is useful for a graceful
// shutdown to make sure all writes/reads are complete before exiting.
//
// Note: I chose to use a composite error here for error handling,
// so that the caller doesn't have to solve one problem in order to
// get to the next (if one exists). So instead of returning on an error
// when we receive it we make our way through all data sources first,
// then generate a composite error with all errors we received along the way,
// append it to the AutocompleteService.Errors list and return it.
//
// With this approach we no longer need a complex management system for in
// place for the Errors slice on our service.
func (a *AutocompleteService) Close() error {
	if a.isClosed {
		return nil
	}
	// Check SnapshotDest DataSource
	var errs []error
	snpErr := a.Config.SnapshotDest.Provider.Close()
	if snpErr != nil {
		errs = append(errs, snpErr)
	}

	for i := range a.Config.DataSources {
		err := a.Config.DataSources[i].Provider.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		compositeErr := fmt.Errorf("autocompleteservice: close: encountered %d errors while closing data sources: %v", len(errs), errs)
		a.Errors = append(a.Errors, compositeErr)
		return compositeErr
	}

	// no need to run GC our service is exiting.
	a.Clear(false)

	a.isClosed = true

	return nil
}

func (a *AutocompleteService) LoadDataSources() error {
	if a.isClosed {
		return fmt.Errorf("autocompleteservice: loaddatasources: service is closed.")
	}

	for _, source := range a.Config.DataSources {
		err := source.Provider.ReadData(source.Filepath, a.store, source.Formatter)
		if err != nil {
			a.Errors = append(a.Errors, err)
			return err
		}
	}
	a.LastUpdated = time.Now().Unix()

	return nil
}

func (a *AutocompleteService) AddSnapshotDest(dest DataSource) {
	a.Config.SnapshotDest = dest
}

func (a *AutocompleteService) CreateSnapshot() error {
	if a.isClosed {
		return fmt.Errorf("autocompleteservice: loaddatasources: service is closed.")
	}
	err := a.Config.SnapshotDest.Provider.DumpData(a.Config.SnapshotDest.Filepath, a.store, a.Config.SnapshotDest.Formatter)
	if err != nil {
		a.Errors = append(a.Errors, err)
	}
	return err
}

func (a *AutocompleteService) RestoreFromSnapshot() error {
	if a.isClosed {
		return fmt.Errorf("autocompleteservice: loaddatasources: service is closed.")
	}
	err := a.Config.SnapshotDest.Provider.ReadData(a.Config.SnapshotDest.Filepath, a.store, a.Config.SnapshotDest.Formatter)
	if err != nil {
		a.Errors = append(a.Errors, err)
		return err
	}
	a.LastUpdated = time.Now().Unix()
	return err
}

func (a *AutocompleteService) LoadDataSource(src DataSource) error {
	if a.isClosed {
		return fmt.Errorf("autocompleteservice: loaddatasources: service is closed.")
	}
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

// Clear will remove all data from the store, in the event you want to start fresh.
// There are two ways we can approach this, the safe way and just set an empty node
// to the root, and just wait for the GC take care of the old one.
//
// Or we could manually trigger a GC cycle. Which is strongly discouraged, but might
// be required in the event of a memory shortage.
//
// You may pass a flag to this function if you wish to manually trigger the GC cycle.
// Please note that running GC manually can:
//
//	Block the caller until the garbage collection is complete.
//	It may also block the entire program.
//	Per the runtime.DC() godocs.
func (a *AutocompleteService) Clear(runGC bool) {
	a.LastUpdated = time.Now().Unix()

	a.store.Clear()
	// TODO: Check to see if just setting the store to nil or creating a new empty store
	// is enough to remove all references to the old data and trigger the GC.

	if runGC {
		runtime.GC()
	}
}

// I am providing different names to these functions to avoid
// implementing the internal interface autocompleter on itself.
// This also provides quick access instead of having to go through
// the store. And gives us room to add more functionality later.
func (a *AutocompleteService) Complete(prefix string) []string {
	if a.isClosed {
		return []string{}
	}
	return a.store.Autocomplete(prefix)
}

func (a *AutocompleteService) Exists(word string) bool {
	if a.isClosed {
		return false
	}
	return a.store.Contains(word)
}

func (a *AutocompleteService) Add(word string) {
	if a.isClosed {
		return
	}
	a.store.Insert(word)
}

func (a *AutocompleteService) GetContents(word string) []string {
	if a.isClosed {
		return []string{}
	}
	return a.store.ListContents()
}
