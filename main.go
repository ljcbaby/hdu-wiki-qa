package main

import (
	"os"

	"github.com/ljcbaby/hdu-wiki-qa/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetOutput(os.Stdout)
	cmd.Execute()
}
