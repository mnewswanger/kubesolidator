package kubeDigests

import (
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"

	"go.mikenewswanger.com/utilities/executil"
	"go.mikenewswanger.com/utilities/filesystem"
	"go.mikenewswanger.com/utilities/slices"
)

// Ordered list of kube object types to loop over
var kubernetesObjectTypes = []kubernetesObjectTypeStruct{
	kubernetesObjectTypeStruct{
		name:  "namespace",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "clusterrole",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "clusterrolebinding",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "role",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "rolebinding",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "serviceaccount",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "persistentvolume",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "persistentvolumeclaim",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "configmap",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "deployment",
		apply: true,
	},
	kubernetesObjectTypeStruct{
		name:  "service",
		apply: true,
	},
}

type kubernetesObjectTypeStruct struct {
	name  string
	apply bool
}

var annotationPrefix = "go.mikenewswanger.com/kubesolidator"

// Apply validates and applies the desired configuration to the cluster
func (kd *KubernetesDigests) Apply(kubectlContext string, dryRun bool) {
	executil.SetVerbosity(verbosity)
	if verbosity > 0 {
		color.White("Validating...")
	}
	if !kd.Validates(verbosity > 0) {
		color.Red("Validation Failed - Exiting")
		os.Exit(1)
	}
	if verbosity > 0 {
		color.Green("Validation Succeeded")
		println()
	}

	// Sort kubeObjects by Kind
	kubeObjectsByKind := make(map[string][]*kubeObject)
	for _, d := range kd.digests {
		kubeObjectsByKind[d.kind] = append(kubeObjectsByKind[d.kind], d)
	}

	for _, kot := range kubernetesObjectTypes {
		if !kot.apply {
			continue
		}
		if verbosity > 0 {
			color.White("Determining deltas for type " + kot.name)
		}
		kubernetesExisting := loadKubernetesObjects(kubectlContext, kot.name)

		objectsInDigest := map[string]bool{}
		objectsToAdd := map[string]*kubeObject{}
		objectsToUpdate := map[string]*kubeObject{}
		objectsToRemove := map[string]bool{}

		for _, o := range kubeObjectsByKind[kot.name] {
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
				if kot.name == "namespace" && (k == "kube-public:kube-public" || k == "kube-system:kube-system") {
					continue
				}
				if kot.name == "service" && k == "default:kubernetes" {
					continue
				}
				objectsToRemove[k] = true
			}
		}

		if dryRun || verbosity > 0 {
			color.White("Changes to apply for: " + kot.name)
			changesMade := false
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
			tempDir, err := ioutil.TempDir(kd.BaseDirectory+"/", ".tmp-")
			if !slices.ContainsString([]string{
				"clusterrole",
				"clusterrolebinding",
				"role",
				"rolebinding",
			}, kot.name) {
				for o := range objectsToRemove {
					deleteKubernetesObject(kubectlContext, kot.name, o)
				}
			}
			for _, o := range objectsToAdd {
				o.apply(tempDir, kubectlContext)
			}
			for _, o := range objectsToUpdate {
				o.apply(tempDir, kubectlContext)
			}
			handleError(err)
			filesystem.RemoveDirectory(tempDir, true)
		}
	}
}

func (ko *kubeObject) apply(tempDir string, kubectlConext string) {
	filename := tempDir + "/" + ko.thumbprint
	yamlData, err := yaml.Marshal(ko.validatedData)
	handleError(err)
	filesystem.WriteFile(filename, yamlData, 0644)
	applyKubernetesObject(kubectlConext, filename)
}

func (ko *kubeObject) addAnnotation(name string, value string) {
	m := ko.validatedData["metadata"].(map[string]interface{})
	a, exists := m["annotations"].(map[string]interface{})
	if !exists {
		a = make(map[string]interface{})
	}
	a[annotationPrefix+"."+name] = value
	m["annotations"] = a
	ko.validatedData["metadata"] = m
}
