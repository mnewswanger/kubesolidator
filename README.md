# Kubesolidator #

Kubsolidator is designed to provide desired state configuration against Kubernetes clusters.

Files should be formatted in the following pattern relative to the file base path:
`/<namespace>/[subfoldersForOrganization]/<objectName>.<objectType>.yml`

It uses `kubectl` to interact with Kubernetes.
