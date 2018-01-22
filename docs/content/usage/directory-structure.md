+++
title = "Directory Structure"
description = "File / Folder structure for digests used by Kubesolidator"
weight = 0
+++

Kubesolidator enforces some directory strucuture and file naming conventions.

### Folder Structure ###

Kubesolidator enforces that the root folders of the directory structure match the namespaces of the objects contained within them.  Inside the top-level subdirectories, subfolders can be used as desired for object organization.

### File Names ###

File names should be `<object-name>.<deployment-type>.yml`.  This makes it very easy to use tools like find to locate objects within a directory structure as well as make it clear looking at the file lists to see the objects maintained in your repository.

### Examples ###

* Namespace (Name: web-apps): `/web-apps/web-apps.namespace.yml`
* Service (Name: example, Namespace: web-apps): `/web-apps/example.service.yml`
* Deployment (Name: example, Namespace: web-apps): `/web-apps/example.deployment.yml`
