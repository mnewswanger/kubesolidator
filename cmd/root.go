package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"go.mikenewswanger.com/kubesolidator/kube-digests"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubesolidator",
	Short: "",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch flags.verbosity {
		case 0:
			logger.Level = logrus.ErrorLevel
			break
		case 1:
			logger.Level = logrus.WarnLevel
			break
		case 2:
			fallthrough
		case 3:
			logger.Level = logrus.InfoLevel
			break
		default:
			logger.Level = logrus.DebugLevel
			break
		}

		kubeDigests.SetLogger(logger)
		kubeDigests.SetVerbosity(uint8(flags.verbosity))

		logger.Debug("Pre-run complete")
	},
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
}

type commandLineFlags struct {
	dryRun                    bool
	kubectlContext            string
	kubernetesDigestDirectory string

	verbosity int
}

var flags = commandLineFlags{}
var logger = logrus.New()
