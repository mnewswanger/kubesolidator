package kubeDigests

import (
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"

	"go.mikenewswanger.com/utilities/filesystem"
)

// Apply validates and applies the desired configuration to the cluster
func (kd *KubernetesDigests) Apply(kubectlContext string, dryRun bool, debug bool, verbosity uint8) {
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

	// Ordered list of kube object kinds to loop over
	for _, kind := range []string{
		"namespace",
		"serviceaccount",
		"persistentvolume",
		"persistentvolumeclaim",
		"configmap",
		"deployment",
		"service",
	} {
		if verbosity > 0 {
			color.White("Determining deltas for type " + kind)
		}
		var kubernetesExisting = loadKubernetesObjects(kubectlContext, kind, debug, verbosity)

		var objectsInDigest = make(map[string]bool)
		var objectsToAdd = make(map[string]*kubeObject)
		var objectsToUpdate = make(map[string]*kubeObject)
		var objectsToRemove = make(map[string]bool)

		for _, o := range kubeObjectsByKind[kind] {
			o.addAnnotation("thumbprint", o.thumbprint)
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
				if kind == "namespace" && (k == "kube-public:kube-public" || k == "kube-system:kube-system") {
					continue
				}
				if kind == "service" && k == "default:kubernetes" {
					continue
				}
				objectsToRemove[k] = true
			}
		}

		if dryRun || verbosity > 0 {
			color.White("Changes to apply for: " + kind)
			var changesMade bool
			if len(objectsToAdd) > 0 {
				color.White("  The following will be added:")
				for item := range objectsToAdd {
					color.Green("    " + item)
				}
				changesMade = true
			}
			if len(objectsToUpdate) > 0 {
				color.White("  The following will be updated:")
				for item := range objectsToUpdate {
					color.Yellow("    " + item)
				}
				changesMade = true
			}
			if len(objectsToRemove) > 0 {
				color.White("  The following will be removed:")
				for item := range objectsToRemove {
					color.Red("    " + item)
				}
				changesMade = true
			}
			if !changesMade {
				color.White("  No changes were made")
			}
			color.White("")
		}

		if !dryRun {
			var tempDir, err = ioutil.TempDir(kd.BaseDirectory+"/", ".tmp-")
			for o := range objectsToRemove {
				deleteKubernetesObject(kubectlContext, kind, o, debug, verbosity)
			}
			for _, o := range objectsToAdd {
				o.apply(tempDir, kubectlContext, debug, verbosity)
			}
			for _, o := range objectsToUpdate {
				o.apply(tempDir, kubectlContext, debug, verbosity)
			}
			handleError(err)
			filesystem.RemoveDirectory(tempDir, true)
		}
	}
}

func (ko *kubeObject) apply(tempDir string, kubectlConext string, debug bool, verbosity uint8) {
	var filename = tempDir + "/" + ko.thumbprint
	var yamlData, err = yaml.Marshal(ko.validatedData)
	handleError(err)
	filesystem.WriteFile(filename, yamlData, 0644)
	applyKubernetesObject(kubectlConext, filename, debug, verbosity)
}

func (ko *kubeObject) addAnnotation(name string, value string) {
	var m = ko.validatedData["metadata"].(map[string]interface{})
	var a, exists = m["annotations"].(map[string]interface{})
	if !exists {
		a = make(map[string]interface{})
	}
	a["kubesolidator."+name] = value
	m["annotations"] = a
	ko.validatedData["metadata"] = m
}
