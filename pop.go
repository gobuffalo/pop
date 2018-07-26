package pop

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

// AvailableDialects lists the available database dialects
var AvailableDialects = []string{"postgres", "mysql", "cockroach"}

// Debug mode, to toggle verbose log traces
var Debug = false

// Color mode, to toggle colored logs
var Color = true
var logger = log.New(os.Stdout, "[POP] ", log.LstdFlags)

// EagerMode type for all eager modes supported in pop.
type EagerMode int8

const (
	eagerModeNil EagerMode = iota
	// EagerDefault is the current implementation, the default
	// behavior of pop. This one introduce N+1 problem and will be used as
	// default value for backward compatibility.
	EagerDefault

	// EagerPreload mode works similar to Preload mode used in ActiveRecord.
	// Avoid N+1 problem by reducing the number of hits to the database but
	// increase memory use to process and link associations to parent.
	EagerPreload

	// EagerInclude This mode works similar to Include mode used in rails ActiveRecord.
	// Use Left Join clauses to load associations.
	EagerInclude
)

// default loading Association Strategy definition.
var loadingAssociationsStrategy = EagerDefault

// SetEagerMode changes overall mode when eager loading.
// this will change the default loading associations strategy for all Eager queries.
// This should be used once, when setting up pop connection.
// func SetEagerMode(eagerMode EagerMode) {
// 	loadingAssociationsStrategy = eagerMode
// }

// Log a formatted string to the logger
var Log = func(s string, args ...interface{}) {
	if Debug {
		if len(args) > 0 {
			xargs := make([]string, len(args))
			for i, a := range args {
				switch a.(type) {
				case string:
					xargs[i] = fmt.Sprintf("%q", a)
				default:
					xargs[i] = fmt.Sprintf("%v", a)
				}
			}
			s = fmt.Sprintf("%s | %s", s, xargs)
		}
		if Color {
			s = color.YellowString(s)
		}
		logger.Println(s)
	}
}

// DialectSupported checks support for the given database dialect
func DialectSupported(d string) bool {
	for _, ad := range AvailableDialects {
		if ad == d {
			return true
		}
	}
	return false
}
