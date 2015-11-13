package pop

import (
	"log"
	"os"
)

var Debug = false
var logger = log.New(os.Stdout, "[POP] ", log.LstdFlags)

var Log = func(s string) {
	if Debug {
		logger.Println(s)
	}
}
