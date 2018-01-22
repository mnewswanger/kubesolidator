+++
title = "Deploying Objects"
description = "Deploy Kubernetes objects to a new or existing cluster"
weight = 10
+++

### Deploying Objects ###

Deploying objects is done via the `kubesolidator apply` commmand.  The command supports the `--dry-run` flag to run the delta operations without modifying any objects.

Kubernetes objects created within the deployment filesystem are simply Kubernetes native YAML objects.  See the [Kubernetes documentation pages](https://kubernetes.io/docs/home/) for details.
