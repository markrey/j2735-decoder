package decoder

import (
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
)

// Logger is used for logging to stdout, thread-safe
var Logger log.Logger

func init() {
	Logger = stdlog.GetFromFlags()
}
