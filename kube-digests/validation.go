package kubeDigests

import (
	"strings"

	"github.com/fatih/color"
)

func (ko *kubeObject) addValidationError(error string) {
	ko.validationErrors = append(ko.validationErrors, error)
}

// Validates validates properties of the Kubernetes objects defined in the digest
func (kd *KubernetesDigests) Validates() bool {
	hasErrors := false
	for _, errors := range kd.Validate() {
		if len(errors) > 0 {
			for _, e := range errors {
				color.Red(e)
			}
			hasErrors = true
		}
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
	if namespace, asserted := m["namespace"].(string); asserted {
		switch ko.kind {
		case "namespace":
			ko.namespace = ko.name
			break
		case "persistentvolume":
			ko.namespace = "_"
			break
		default:
			ko.namespace = namespace
		}

		if !strings.HasPrefix(ko.relativePath, "/"+namespace+"/") {
			ko.addValidationError("Should exist in proper namespace folder (/" + namespace + "/)")
		}
	} else {
		ko.addValidationError("Missing 'namespace' property in metadata")
	}

	ko.validatedData["metadata"] = m
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
