package kubeDigests

import (
	"strings"

	"github.com/ghodss/yaml"

	"go.mikenewswanger.com/utilities/filesystem"
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
	contents, err := filesystem.LoadFileBytes(ko.absolutePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(contents, &ko.rawData)
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
	directoryContents, err := filesystem.GetDirectoryContents(kd.BaseDirectory + subfolder)

	if err != nil {
		panic(err)
	}
	for _, item := range directoryContents {
		if strings.HasPrefix(item, ".") {
			continue
		}
		absolutePath = kd.BaseDirectory + subfolder + item
		if filesystem.IsDirectory(absolutePath) {
			kd.loadDigestsInFolder(subfolder + item + "/")
		} else {
			var checksum string
			checksum, err = filesystem.GetFileSHA256Checksum(absolutePath)
			if err != nil {
				panic(err)
			}
			var obj = kubeObject{
				absolutePath: absolutePath,
				relativePath: subfolder + item,
				thumbprint:   checksum,
			}
			obj.loadDataFromFile()
			kd.digests = append(kd.digests, &obj)
		}
	}
}
