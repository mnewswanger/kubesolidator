package kubeDigests

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func handleError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
