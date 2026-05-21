package cmd

import (
	"io"
	"os"
	"os/exec"
)

type HelmRunner interface {
	Run(args []string, stdout, stderr io.Writer) error
}

type execRunner struct{}

func (execRunner) Run(args []string, stdout, stderr io.Writer) error {
	helmBin := os.Getenv("HELM_BIN")
	if helmBin == "" {
		helmBin = "helm"
	}

	cmd := exec.Command(helmBin, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

var runner HelmRunner = execRunner{}

func setRunner(r HelmRunner) {
	runner = r
}
