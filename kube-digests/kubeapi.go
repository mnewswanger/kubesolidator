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
	items []struct {
		metadata struct {
			name        string
			namespace   string
			annotations map[string]string
		}
	}
}

func loadKubernetesObjects(kubernetesAPIServer string, apiEndpointType string) map[string]string {
	var kubeObjectList = make(map[string]string)

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

	for _ = range kubernetesObjects.items {
		panic("looping")
	}

	return kubeObjectList
}

// {
//   "kind": "DeploymentList",
//   "apiVersion": "extensions/v1beta1",
//   "metadata": {
//     "selfLink": "/apis/extensions/v1beta1/deployments",
//     "resourceVersion": "8223147"
//   },
//   "items": [
//     {
//       "metadata": {
//         "name": "packages-webui",
//         "namespace": "build-servers",
//         "selfLink": "/apis/extensions/v1beta1/namespaces/build-servers/deployments/packages-webui",
//         "uid": "faffe29f-0b26-11e7-8a3c-005056a22086",
//         "resourceVersion": "8219879",
//         "generation": 2,
//         "creationTimestamp": "2017-03-17T15:32:50Z",
//         "labels": {
//           "app": "packages",
//           "role": "web-ui"
//         },
//         "annotations": {
//           "deployment.kubernetes.io/revision": "2"
//         }
//       },
//       "spec": {
//         "replicas": 1,
//         "selector": {
//           "matchLabels": {
//             "app": "packages",
//             "role": "web-ui"
//           }
//         },
//         "template": {
//           "metadata": {
//             "creationTimestamp": null,
//             "labels": {
//               "app": "packages",
//               "role": "web-ui"
//             }
//           },
//           "spec": {
//             "volumes": [
//               {
//                 "name": "builds",
//                 "persistentVolumeClaim": {
//                   "claimName": "build-server-binaries"
//                 }
//               }
//             ],
//             "containers": [
//               {
//                 "name": "nginx",
//                 "image": "docker-registry.home.mikenewswanger.com/build-servers/web-interface:develop",
//                 "resources": {},
//                 "volumeMounts": [
//                   {
//                     "name": "builds",
//                     "readOnly": true,
//                     "mountPath": "/var/www/html"
//                   }
//                 ],
//                 "terminationMessagePath": "/dev/termination-log",
//                 "imagePullPolicy": "Always"
//               }
//             ],
//             "restartPolicy": "Always",
//             "terminationGracePeriodSeconds": 30,
//             "dnsPolicy": "ClusterFirst",
//             "securityContext": {}
//           }
//         },
//         "strategy": {
//           "type": "RollingUpdate",
//           "rollingUpdate": {
//             "maxUnavailable": 1,
//             "maxSurge": 1
//           }
//         }
//       },
//       "status": {
//         "observedGeneration": 2,
//         "replicas": 1,
//         "updatedReplicas": 1,
//         "availableReplicas": 1,
//         "conditions": [
//           {
//             "type": "Available",
//             "status": "True",
//             "lastUpdateTime": "2017-03-17T15:32:51Z",
//             "lastTransitionTime": "2017-03-17T15:32:51Z",
//             "reason": "MinimumReplicasAvailable",
//             "message": "Deployment has minimum availability."
//           }
//         ]
//       }
//     },
