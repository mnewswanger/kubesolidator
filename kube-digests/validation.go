package kubeDigests

import (
	"strings"
)

func (ko *kubeObject) addValidationError(error string) {
	ko.validationErrors = append(ko.validationErrors, error)
}

// Validates validates properties of the Kubernetes objects defined in the digest
func (kd *KubernetesDigests) Validates() bool {
	for _, errors := range kd.Validate() {
		if len(errors) > 0 {
			return false
		}
	}
	return true
}

// Validate validates properties of the Kubernetes objects defined in the digest
func (kd *KubernetesDigests) Validate() map[string][]string {
	kd.loadDigests()

	var validationErrors = make(map[string][]string)

	for _, ko := range kd.digests {
		d := ko.data.(map[string]interface{})
		if kind, asserted := d["kind"].(string); asserted {
			ko.validateBaseMetadata(kind, d["metadata"])
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
			println(ko.absolutePath + ": " + kind)
		} else {
			ko.addValidationError("Kind not specified on object")
		}

		if len(ko.validationErrors) > 0 {
			validationErrors[ko.relativePath] = ko.validationErrors
		}
	}

	return validationErrors
}

func (ko *kubeObject) validateBaseMetadata(kind string, metadata interface{}) {
	m := metadata.(map[string]interface{})

	// Validate filename
	if name, asserted := m["name"].(string); asserted {
		var filenameShouldBe = name + "." + strings.ToLower(kind) + ".yml"
		if !strings.HasSuffix(ko.relativePath, filenameShouldBe) {
			ko.addValidationError("Imporoperly named file (should be " + filenameShouldBe + ")")
		}
	} else {
		ko.addValidationError("Missing 'name' property in metadata")
	}

	// Validate namespace folder structure
	if namespace, asserted := m["namespace"].(string); asserted {
		if !strings.HasPrefix(ko.relativePath, "/"+namespace+"/") {
			ko.addValidationError("Should exist in proper namespace folder (/" + namespace + "/)")
		}
	} else {
		ko.addValidationError("Missing 'namespace' property in metadata")
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
