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
	var fs = filesystem.Filesystem{}
	var contents, err = fs.LoadFileIfExists(ko.absolutePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(contents), &ko.rawData)
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
	var fs = filesystem.Filesystem{}
	var absolutePath string
	var directoryContents, err = fs.GetDirectoryContents(kd.BaseDirectory + subfolder)

	if err != nil {
		panic(err)
	}
	for _, item := range directoryContents {
		if strings.HasPrefix(item, ".") {
			continue
		}
		absolutePath = kd.BaseDirectory + subfolder + item
		if fs.IsDirectory(absolutePath) {
			kd.loadDigestsInFolder(subfolder + item + "/")
		} else {
			var checksum string
			checksum, err = fs.GetFileSHA256Checksum(absolutePath)
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
