package kubeDigests

import (
	"os"

	"github.com/fatih/color"
)

// Apply validates and applies the desired configuration to the cluster
func (kd *KubernetesDigests) Apply(kubernetesAPIServer string, dryRun bool) {
	if !kd.Validates() {
		color.Red("Validation Failed - Exiting")
		os.Exit(1)
	}
}
