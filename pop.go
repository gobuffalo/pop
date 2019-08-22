package pop

import "strings"

// AvailableDialects lists the available database dialects
var AvailableDialects []string

var dialectSynonyms = make(map[string]string)

// map of dialect specific url parsers
var urlParser = make(map[string]func(*ConnectionDetails) error)

// map of dialect specific connection details finalizers
var finalizer = make(map[string]func(*ConnectionDetails))

// map of connection creators
var newConnection = make(map[string]func(*ConnectionDetails) (dialect, error))

// DialectSupported checks support for the given database dialect
func DialectSupported(d string) bool {
	for _, ad := range AvailableDialects {
		if ad == d {
			return true
		}
	}
	return false
}

func normalizeSynonyms(dialect string) string {
	d := strings.ToLower(dialect)
	if syn, ok := dialectSynonyms[d]; ok {
		d = syn
	}
	return d
}
