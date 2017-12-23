package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"go.mikenewswanger.com/kubesolidator/kube-digests"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		kd := kubeDigests.KubernetesDigests{
			BaseDirectory: flags.kubernetesDigestDirectory,
		}
		if kd.Validates(true) {
			color.Green("Digests validated successfully")
			println()
		} else {
			color.Red("Digest validation failed successfully")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
