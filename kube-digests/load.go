package kubeDigests

import (
	"strings"

	"github.com/ghodss/yaml"

	"gitlab.home.mikenewswanger.com/golang/filesystem"
)

// KubernetesDigests provides interaaction to the digest repository contents
type KubernetesDigests struct {
	BaseDirectory string
	digestsLoaded bool
	digests       []*kubeObject
}

type kubeObject struct {
	// File Data
	absolutePath string
	relativePath string
	thumbprint   string

	// Derived Data
	kind          string
	name          string
	namespace     string
	validatedData map[string]interface{}

	// Imported Data
	rawData interface{}

	// Validation
	validationErrors []string
}

func (ko *kubeObject) loadDataFromFile() {
	var err = yaml.Unmarshal([]byte(filesystem.LoadFileIfExists(ko.absolutePath)), &ko.rawData)
	if err != nil {
		panic(err)
	}
}

func (kd *KubernetesDigests) loadDigests() {
	if !kd.digestsLoaded {
		kd.loadDigestsInFolder("/")
		kd.digestsLoaded = true
	}
}

func (kd *KubernetesDigests) loadDigestsInFolder(subfolder string) {
	var absolutePath string
	for _, item := range filesystem.GetDirectoryContents(kd.BaseDirectory + subfolder) {
		if strings.HasPrefix(item, ".") {
			continue
		}
		absolutePath = kd.BaseDirectory + subfolder + item
		if filesystem.IsDirectory(absolutePath) {
			kd.loadDigestsInFolder(subfolder + item + "/")
		} else {
			var obj = kubeObject{
				absolutePath: absolutePath,
				relativePath: subfolder + item,
				thumbprint:   filesystem.GetFileSHA256Checksum(absolutePath),
			}
			obj.loadDataFromFile()
			kd.digests = append(kd.digests, &obj)
		}
	}
}
