package pop

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/fatih/color"
	"github.com/gobuffalo/pop/v6/logging"
)

// Debug mode, to toggle verbose log traces
var Debug = false

// Color mode, to toggle colored logs
var Color = true

// SetLogger overrides the default logger.
func SetLogger(logger func(level logging.Level, s string, args ...interface{})) {
	log = logger
}

var defaultStdLogger = stdlog.New(os.Stderr, "[POP] ", stdlog.LstdFlags)

var log = func(lvl logging.Level, s string, args ...interface{}) {
	if !Debug && lvl <= logging.Debug {
		return
	}
	if lvl == logging.SQL {
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
			s = fmt.Sprintf("%s - %s | %s", lvl, s, xargs)
		} else {
			s = fmt.Sprintf("%s - %s", lvl, s)
		}
	} else {
		s = fmt.Sprintf(s, args...)
		s = fmt.Sprintf("%s - %s", lvl, s)
	}
	if Color {
		s = color.YellowString(s)
	}
	defaultStdLogger.Println(s)
}
