package main

import (
	"os"

	"github.com/postfinance/kuota-calc/cmd"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

//nolint:gochecknoglobals
var (
	// set by goreleaser on build
	version, date, commit string = "master", "?", "?"
)

const (
	binaryName = "kuota-calc"
)

func main() {
	flags := pflag.NewFlagSet(binaryName, pflag.ExitOnError)
	pflag.CommandLine = flags

	v := cmd.Version{
		Version: version,
		Date:    date,
		Commit:  commit,
	}

	root := cmd.NewKuotaCalcCmd(&v, genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
