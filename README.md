# Autocomplete
## 🚧  Under construction. - Coming soon  🚧 
A project to explore Trie's, Ternary Search Trees, dfs, and more!

## Overview
Autocomplete provides an in-memory autocompletion engine.
It is designed to be used in conjunction with a data store to provide
autocompletion suggestions for a given prefix. This is useful for things
like search bars, command line interfaces, and more.


##### WIP
- [X] Create some sort of saveSnapshot locally or to buckets even??? 
- [X] Enable data sources that are responsible for populating the keyword list.
- [X] Add coniguration option to choose the type of tree to use.
- [ ] Complete Default Formatter tests
- [ ] Complete KeywordList Formatter tests
- [ ] Complete LocalFileProvider tests
- [ ] Complete GoogleStorageBucketProvider tests
- [ ] Complete GithubProvider tests
- [ ] Complete AutoCompleteService tests
- [ ] Add any missing `New` methods.
- [ ] Work on complete sample service
- [ ] Benchmarks w/ Examples
- [ ] Cleanup unused properties/settings.
- [ ] Setup cmd/autocompleter cli tool to run an AutoCompleteService 



##### FUTURE
- [ ] Profile to try and improve memory, performance and GC time.
- [ ] Sharding Trees across mutexes
- [ ] Export with binary instead of plain text
- [ ] Extend by adding an additional Data structure 
- [ ] What would parallel computing look like?  
