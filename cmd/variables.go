package cmd

type commandLineFlags struct {
	dryRun                    bool
	kubectlContext            string
	kubernetesDigestDirectory string

	debug     bool
	verbosity int
}

var flags = commandLineFlags{}
