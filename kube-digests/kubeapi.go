package kubeDigests

import (
	"strings"

	"github.com/ghodss/yaml"

	"go.mikenewswanger.com/utilities/executil"
)

type kubernetesObjectsStruct struct {
	Items []struct {
		Metadata struct {
			Name        string
			Namespace   string
			Annotations map[string]string
		}
	}
}

func applyKubernetesObject(kubectlContext string, file string, debug bool, verbosity uint8) {
	var args = []string{}
	if kubectlContext != "" {
		args = append(args, "--context", kubectlContext)
	}
	args = append(args, "apply", "-f", file)
	var c = executil.Command{
		Name:       "Apply kubernetes object (" + file + ")",
		Executable: "kubectl",
		Arguments:  args,
	}
	if err := c.Run(); err != nil {
		panic(err)
	}
}

func deleteKubernetesObject(kubectlContext string, kind string, item string, debug bool, verbosity uint8) {
	var args = []string{}

	if kubectlContext != "" {
		args = append(args, "--context", kubectlContext)
	}

	args = append(args, "delete", kind)

	// [0] = namespace; [1] = name
	var kubeObjectParts = strings.Split(item, ":")
	if kubeObjectParts[0] != "_" {
		args = append(args, "--namespace", kubeObjectParts[0])
	}

	args = append(args, kubeObjectParts[1])

	var c = executil.Command{
		Name:       "Remove kubernetes object: " + kind + " - " + item,
		Executable: "kubectl",
		Arguments:  args,
	}
	if e := c.Run(); e != nil {
		panic(e)
	}
}

// loadKubernetesObjects by type
func loadKubernetesObjects(kubectlContext string, kind string, debug bool, verbosity uint8) map[string]string {
	var err error
	var args = []string{}
	if kubectlContext != "" {
		args = []string{"--context", kubectlContext}
	}
	args = append(args, "get", "--all-namespaces", "-o", "yaml", kind)

	// Get the existing objects from kubectl
	var c = executil.Command{
		Name:       "Load kubernetes objects (kind: " + kind + ")",
		Executable: "kubectl",
		Arguments:  args,
	}

	err = c.Run()
	handleError(err)

	var output = c.GetStdout()
	var kubernetesObjects = kubernetesObjectsStruct{}
	err = yaml.Unmarshal([]byte(output), &kubernetesObjects)
	handleError(err)

	var kubeObjectList = make(map[string]string)
	for _, item := range kubernetesObjects.Items {
		switch kind {
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
