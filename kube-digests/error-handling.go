package kubeDigests

import (
	"github.com/sirupsen/logrus"

	"go.mikenewswanger.com/utilities/executil"
)

var verbosity = uint8(0)
var logger = logrus.New()

// SetLogger sets the logger for the package
func SetLogger(l *logrus.Logger) {
	logger = l
	executil.SetLogger(l)
}

// SetVerbosity sets the verbosity level for the package
func SetVerbosity(v uint8) {
	verbosity = v
	executil.SetVerbosity(v)
}

func handleError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
