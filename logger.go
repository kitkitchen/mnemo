package mnemo

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "Mnemo: ",
})
