package cmd

import (
	"github.com/spf13/cobra"

	"go.mikenewswanger.com/kubesolidator/kube-digests"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var kd = kubeDigests.KubernetesDigests{
			BaseDirectory: flags.kubernetesDigestDirectory,
		}
		kd.Apply(flags.kubectlContext, flags.dryRun)
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVarP(&flags.dryRun, "dry-run", "", false, "Perform a no-op")
	applyCmd.Flags().StringVarP(&flags.kubectlContext, "kubectl-context", "k", "", "Kubectl context to operate against")
}
