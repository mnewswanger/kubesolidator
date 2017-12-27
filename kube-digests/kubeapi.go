package kubeDigests

import (
	"bufio"
	"strings"

	"github.com/fatih/color"
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

func applyKubernetesObject(kubectlContext string, file string) {
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

func deleteKubernetesObject(kubectlContext string, kind string, item string) {
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
func loadKubernetesObjects(kubectlContext string, kind string) map[string]string {
	args := []string{}
	if kubectlContext != "" {
		args = []string{"--context", kubectlContext}
	}
	args = append(
		args,
		"get",
		"--all-namespaces",
		"--output",
		"template",
		kind,
		"--template",
		"{{range $k, $i := .items }}{{$i.metadata.namespace}}:{{$i.metadata.name}} {{if $i.metadata.annotations}}{{index $i.metadata.annotations \"kubesolidator.thumbprint\"}}{{end}}\n{{end}}",
	)

	// Get the existing objects from kubectl
	c := executil.Command{
		Name:       "Load kubernetes objects (kind: " + kind + ")",
		Executable: "kubectl",
		Arguments:  args,
	}

	err := c.Run()
	handleError(err)

	scanner := bufio.NewScanner(strings.NewReader(c.GetStdout()))
	kubeObjectList := make(map[string]string)
	if verbosity > 1 {
		color.Blue("Objects of type already in cluster: " + kind)
	}
	for scanner.Scan() {
		line := scanner.Text()
		if verbosity > 1 {
			color.Blue("  " + line)
		}
		split := strings.Fields(line)
		thumbprint := ""
		if len(split) > 1 {
			thumbprint = split[1]
		}
		kubeObjectList[split[0]] = thumbprint
	}
	return kubeObjectList
}
