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
		var kd = kubeDigests.KubernetesDigests{
			BaseDirectory: flags.kubernetesDigestDirectory,
		}
		var hasErrors = false
		for path, errors := range kd.Validate() {
			hasErrors = true
			color.Red(path)
			for _, error := range errors {
				color.Red("  " + error)
			}
			println()
		}
		if !hasErrors {
			color.Green("Digests validated successfully")
			println()
		} else {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
