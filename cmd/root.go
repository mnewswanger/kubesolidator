package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubesolidator",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&flags.kubernetesDigestDirectory, "kubernetes-digest-directory", "d", "", "Directory containing Kubernetes digests")
	RootCmd.PersistentFlags().CountVarP(&flags.verbosity, "verbosity", "v", "Output verbosity")
	RootCmd.PersistentFlags().BoolVarP(&flags.debug, "debug", "", false, "Debug level output")
}
