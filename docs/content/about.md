+++
title = "About"
description = "Desired State Deployment for Kubernetes"
+++

Kubesolidator is designed to help codify your infrastructure.  It provides a fast, consistent way to deploy using a gitops workflow.

### Desired State ###

Kubesolidator compares what's in your filesystem to what's in your cluster and applies a delta operation to ensure consistency--both on creation and removal--of your cluster deployments.  It performs these deltas efficiently and uses kubernetes native tooling to do all of its work--ensuring high levels of compatibility with different Kubernetes versions.

### Flexible ###

Kubesolidator can quickly apply deltas to existing clusters or consistently deploy to new clusters from scratch.  This allows teams to manage multiple clusters with very little management overhead.

Once deployed, Kubesolidator tooling provides ability to manipulate the deployments it's made.  This includes the ability to update images and manipulate scaling for integration with CI/CD pipelines.

### Fast ###

By doing comparison operations in bulk, performance can be maintained even with very large clusters.  Maintaining performance at scale is a core objective of the project so that CI/CD systems can run with little latency.

<div id="action-buttons">
    <a class="button primary big" href="{{< ref "getting-started.md" >}}">Get Started</a>
    <a class="button outline big" href="{{< ref "downloads.md" >}}" >Download</a>
</div>
