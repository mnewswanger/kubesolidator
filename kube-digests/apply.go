package kubeDigests

import (
	"os"

	"github.com/fatih/color"
)

// Apply validates and applies the desired configuration to the cluster
func (kd *KubernetesDigests) Apply(kubernetesAPIServer string, dryRun bool, debug bool, verbosity uint8) {
	if verbosity > 0 {
		color.White("Validating...")
	}
	if !kd.Validates() {
		color.Red("Validation Failed - Exiting")
		os.Exit(1)
	}
	if verbosity > 0 {
		color.Green("Validation Succeeded")
		println()
	}

	// Sort kubeObjects by Kind
	var kubeObjectsByKind = make(map[string][]*kubeObject)
	for _, d := range kd.digests {
		kubeObjectsByKind[d.kind] = append(kubeObjectsByKind[d.kind], d)
	}

	for _, kind := range []string{
		"namespace",
		"persistentvolume",
		"persistentvolumeclaim",
		"configmap",
		"deployment",
		"service",
	} {
		if verbosity > 0 {
			color.White("Determining deltas for type " + kind)
		}
		var kubernetesExisting = loadKubernetesObjects(kubernetesAPIServer, kind)

		var objectsInDigest = make(map[string]bool)
		var objectsToAdd = make(map[string]*kubeObject)
		var objectsToUpdate = make(map[string]*kubeObject)
		var objectsToRemove = make(map[string]bool)

		for _, o := range kubeObjectsByKind[kind] {
			objectsInDigest[o.namespace+":"+o.name] = true
			// Add anything in digests that doesn't already exist
			if _, exists := kubernetesExisting[o.namespace+":"+o.name]; !exists {
				objectsToAdd[o.namespace+":"+o.name] = o
				continue
			}

			// Update anything in digests that exist but don't match existing
			if o.thumbprint != kubernetesExisting[o.namespace+":"+o.name] {
				objectsToUpdate[o.namespace+":"+o.name] = o
			}
		}

		// Remove anything that exists but isn't in digests
		for k := range kubernetesExisting {
			if !objectsInDigest[k] {
				objectsToRemove[k] = true
			}
		}

		if dryRun || verbosity > 0 {
			if len(objectsToAdd) > 0 {
				color.White("The following will be added:")
				for item := range objectsToAdd {
					color.Green("  " + item)
				}
			}
			if len(objectsToUpdate) > 0 {
				color.White("The following will be updated:")
				for item := range objectsToUpdate {
					color.Yellow("  " + item)
				}
			}
			if len(objectsToRemove) > 0 {
				color.White("The following will be removed:")
				for item := range objectsToRemove {
					color.Red("  " + item)
				}
			}
		}
	}
}
