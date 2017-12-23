package kubeDigests

import (
	"strings"

	"go.mikenewswanger.com/utilities/slices"

	"github.com/fatih/color"
)

func (ko *kubeObject) addValidationError(error string) {
	ko.validationErrors = append(ko.validationErrors, error)
}

// Validates validates properties of the Kubernetes objects defined in the digest
func (kd *KubernetesDigests) Validates(printErrors bool) bool {
	hasErrors := false
	for path, errors := range kd.Validate() {
		hasErrors = true
		color.Red(path)
		for _, error := range errors {
			color.Red("  " + error)
		}
		println()
	}
	return !hasErrors
}

// Validate validates properties of the Kubernetes objects defined in the digest
func (kd *KubernetesDigests) Validate() map[string][]string {
	kd.loadDigests()

	var validationErrors = make(map[string][]string)

	for _, ko := range kd.digests {
		ko.validatedData = ko.rawData.(map[string]interface{})
		if kind, asserted := ko.validatedData["kind"].(string); asserted {
			ko.validateBaseMetadata()
			switch kind {
			case "ConfigMap":
				ko.validateConfigMap()
				break
			case "Deployment":
				ko.validateDeployment()
				break
			case "Namespace":
				ko.validateNamespace()
				break
			case "PersistentVolume":
				ko.validatePersistentVolume()
				break
			case "PersistentVolumeClaim":
				ko.validatePersistentVolumeClaim()
				break
			case "Service":
				ko.validateService()
				break
			case "ServiceAccount":
				ko.validateService()
				break
			default:
				ko.addValidationError("Unsupported object type: " + kind)
			}
		} else {
			ko.addValidationError("Kind not specified on object")
		}

		if len(ko.validationErrors) > 0 {
			validationErrors[ko.relativePath] = ko.validationErrors
		}
	}

	return validationErrors
}

func (ko *kubeObject) validateBaseMetadata() {
	ko.kind = strings.ToLower(ko.validatedData["kind"].(string))
	m := ko.validatedData["metadata"].(map[string]interface{})

	// Validate filename
	if name, asserted := m["name"].(string); asserted {
		ko.name = name
		var filenameShouldBe = name + "." + ko.kind + ".yml"
		if !strings.HasSuffix(ko.relativePath, filenameShouldBe) {
			ko.addValidationError("Imporoperly named file (should be " + filenameShouldBe + ")")
		}
	} else {
		ko.addValidationError("Missing 'name' property in metadata")
	}

	// Validate namespace folder structure
	namespace, namespaceSpecified := m["namespace"].(string)
	if slices.ContainsString([]string{
		"namespace",
		"persistentvolume",
	}, ko.kind) {
		if namespaceSpecified {
			ko.addValidationError("Contains 'namespace' property and is a global object")
		}
	} else {
		if namespaceSpecified {
			ko.namespace = namespace
		} else {
			ko.addValidationError("Missing 'namespace' property in metadata")
		}
	}

	if len(ko.validationErrors) == 0 {
		ko.validatedData["metadata"] = m
	}
}

func (ko *kubeObject) validateConfigMap() {
}

func (ko *kubeObject) validateDeployment() {
}

func (ko *kubeObject) validateNamespace() {
}

func (ko *kubeObject) validatePersistentVolume() {
}

func (ko *kubeObject) validatePersistentVolumeClaim() {
}

func (ko *kubeObject) validateService() {
}
