+++
title = "Getting Started"
description = "Quick Start Guide"
+++

### Installation ###

To install Kubesolidator, either [download a prebuilt binary]({{< relref "downloads.md" >}}) or build it yourself:

```
go get go.mikenewswanger.com/kubesolidator
go install go.mikenewswanger.com/kubesolidator
```

### Deploying ###

First, [create your deployment file structure and objects]({{< relref "usage/directory-structure.md" >}}) (digest directory), then deploy:

```
# validate the configuration
kubesolidator validate --kubernetes-digest-directory </path/to/digest-directory>
# see what changes will be applied
kubesolidator apply --kubernetes-digest-directory </path/to/digest-directory> --dry-run
# apply the changes
kubesolidator apply --kubernetes-digest-directory </path/to/digest-directory>
```

<div id="action-buttons">
    <a class="button primary big" href="{{< ref "usage/_index.md" >}}">Advanced Usage</a>
</div>
