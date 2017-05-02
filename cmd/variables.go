package cmd

type commandLineFlags struct {
	dryRun                    bool
	kubernetesAPIServer       string
	kubernetesDigestDirectory string

	debug     bool
	verbosity int
}

var flags = commandLineFlags{}
