// Package cmd provides the kuota-calc command.
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"runtime"
	"text/tabwriter"

	"github.com/postfinance/kuota-calc/internal/calc"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	kuotaCalcExample = `    # provide a simple/complex deployment by piping it to kuota-calc (used as kubectl plugin)
    cat deployment.yaml | kubectl %[1]s

    # do the same, calling the binary directly with detailed output
    cat deployment.yaml | %[1]s --detailed`
)

// KuotaCalcOpts holds all command options.
type KuotaCalcOpts struct {
	genericclioptions.IOStreams

	// flags
	debug    bool
	detailed bool
	version  bool
	// files    []string

	versionInfo *Version
}

// NewKuotaCalcCmd returns a coba command wrapping KuotaCalcOps
func NewKuotaCalcCmd(version *Version, streams genericclioptions.IOStreams) *cobra.Command {
	opts := KuotaCalcOpts{
		IOStreams:   streams,
		versionInfo: version,
	}

	cmd := &cobra.Command{
		Use:          "kuota-calc",
		Short:        "Calculate the resource quota needs of your deployment(s).",
		Example:      fmt.Sprintf(kuotaCalcExample, "kuota-calc"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if opts.version {
				return opts.printVersion()
			}

			return opts.run()
		},
	}

	cmd.Flags().BoolVar(&opts.debug, "debug", false, "enable debug logging")
	cmd.Flags().BoolVar(&opts.detailed, "detailed", false, "enable detailed output")
	cmd.Flags().BoolVar(&opts.version, "version", false, "print version and exit")

	return cmd
}

func (opts *KuotaCalcOpts) printVersion() error {
	fmt.Fprintf(opts.Out, "version %s (revision: %s)\n\tbuild date: %s\n\tgo version: %s\n",
		opts.versionInfo.Version,
		opts.versionInfo.Commit,
		opts.versionInfo.Date,
		runtime.Version(),
	)

	return nil
}

func (opts *KuotaCalcOpts) run() error {
	var (
		summary []*calc.ResourceUsage
	)

	yamlReader := yaml.NewYAMLReader(bufio.NewReader(opts.In))

	for {
		data, err := yamlReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("reading input: %w", err)
		}

		usage, err := calc.ResourceQuotaFromYaml(data)
		if err != nil {
			if errors.Is(err, calc.ErrResourceNotSupported) {
				if opts.debug {
					fmt.Fprintf(opts.Out, "DEBUG: %s\n", err)
				}
				continue
			}

			return err
		}

		summary = append(summary, usage)
	}

	if opts.detailed {
		opts.printDetailed(summary)
	} else {
		opts.printSummary(summary)
	}

	return nil
}

func (opts *KuotaCalcOpts) printDetailed(usage []*calc.ResourceUsage) {
	w := tabwriter.NewWriter(opts.Out, 0, 0, 4, ' ', tabwriter.TabIndent)

	fmt.Fprintf(w, "Version\tKind\tName\tReplicas\tStrategy\tMaxReplicas\tCPU\tMemory\t\n")

	for _, u := range usage {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%d\t%s\t%s\t\n",
			u.Details.Version,
			u.Details.Kind,
			u.Details.Name,
			u.Details.Replicas,
			u.Details.Strategy,
			u.Details.MaxReplicas,
			u.CPU,
			u.Memory,
		)
	}

	w.Flush()

	fmt.Fprintf(opts.Out, "\nTotal\n")

	opts.printSummary(usage)
}

func (opts *KuotaCalcOpts) printSummary(usage []*calc.ResourceUsage) {
	var (
		cpuUsage    resource.Quantity
		memoryUsage resource.Quantity
	)

	for _, u := range usage {
		cpuUsage.Add(*u.CPU)
		memoryUsage.Add(*u.Memory)
	}

	fmt.Fprintf(opts.Out, "CPU: %s\nMemory: %s\n",
		cpuUsage.String(),
		memoryUsage.String(),
	)
}
