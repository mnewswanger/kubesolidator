package kubeDigests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var apiEndpoints = map[string]string{
	"configmap":             "/api/v1/configmaps",
	"deployment":            "/apis/extensions/v1beta1/deployments",
	"namespace":             "/api/v1/namespaces",
	"persistentvolume":      "/api/v1/persistentvolumes",
	"persistentvolumeclaim": "/api/v1/persistentvolumeclaims",
	"service":               "/api/v1/services",
}

type kubernetesObjectsStruct struct {
	Items []struct {
		Metadata struct {
			Name        string
			Namespace   string
			Annotations map[string]string
		}
	}
}

// loadKubernetesObjects by type
func loadKubernetesObjects(kubernetesAPIServer string, apiEndpointType string) map[string]string {

	// Get the existing objects from the KubernetesAPI
	var resp, err = http.Get(kubernetesAPIServer + apiEndpoints[apiEndpointType])
	handleError(err)
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	var kubernetesObjects = kubernetesObjectsStruct{}
	err = json.Unmarshal(body, &kubernetesObjects)
	handleError(err)

	var kubeObjectList = make(map[string]string)
	for _, item := range kubernetesObjects.Items {
		switch apiEndpointType {
		case "namespace":
			item.Metadata.Namespace = item.Metadata.Name
			break
		case "persistentvolume":
			item.Metadata.Namespace = "_"
			break
		}
		kubeObjectList[item.Metadata.Namespace+":"+item.Metadata.Name] = item.Metadata.Annotations["kubesolidator.thumbprint"]
	}
	return kubeObjectList
}
