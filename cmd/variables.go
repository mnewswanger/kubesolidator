package cmd

import (
	"github.com/sirupsen/logrus"
)

type commandLineFlags struct {
	dryRun                    bool
	kubectlContext            string
	kubernetesDigestDirectory string

	debug     bool
	verbosity int
}

var flags = commandLineFlags{}
var logger = logrus.New()
