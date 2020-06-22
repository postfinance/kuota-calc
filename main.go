package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"text/tabwriter"

	"github.com/postfinance/kuota-calc/internal/calc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/deprecated/scheme"
)

//nolint:gochecknoglobals
var (
	// set by goreleaser on build
	version, date, commit string = "master", "?", "?"
	// name of the binary
	binaryName = filepath.Base(os.Args[0])
)

func main() {
	var (
		// default to info loglevel
		logLevel zapcore.Level = zap.InfoLevel

		// flags
		debug    = flag.Bool("debug", false, "Enable debug logging")
		version  = flag.Bool("version", false, "Print version and exit")
		detailed = flag.Bool("detailed", false, "Print detailed output per k8s resource")
	)

	flag.Parse()

	if *debug {
		logLevel = zap.DebugLevel
	}

	log := setupZap(logLevel)

	if *version {
		printVersion()

		return
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// stdin must contain input
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		flag.Usage()

		return
	}

	var (
		summary []*calc.ResourceUsage
	)

	r := yaml.NewYAMLReader(bufio.NewReader(os.Stdin))

	for {
		data, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode

		object, gvk, err := decode(data, nil, nil)
		if err != nil {
			log.Fatalf("decoding: %s", err)
			continue
		}

		switch obj := object.(type) {
		case *appsv1.Deployment:
			usage, err := calc.Deployment(*obj)
			if err != nil {
				log.Errorf("calculating deployment resource usage: %s", err)
				continue
			}

			summary = append(summary, usage)
		case *appsv1.StatefulSet:
			usage, err := calc.StatefulSet(*obj)
			if err != nil {
				log.Errorf("calculating statefulset resource usage: %s", err)
				continue
			}

			summary = append(summary, usage)
		default:
			log.Debugf("ignoring %s", gvk)
			continue
		}
	}

	if *detailed {
		printDetailed(summary)
	} else {
		printSimple(summary)
	}
}

func printDetailed(usage []*calc.ResourceUsage) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)

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

	fmt.Printf("\nTotal\n")

	printSimple(usage)
}

func printSimple(usage []*calc.ResourceUsage) {
	var (
		cpuUsage    resource.Quantity
		memoryUsage resource.Quantity
	)

	for _, u := range usage {
		cpuUsage.Add(*u.CPU)
		memoryUsage.Add(*u.Memory)
	}

	fmt.Printf("CPU: %s\nMemory: %s\n",
		cpuUsage.String(),
		memoryUsage.String(),
	)
}

func setupZap(level zapcore.Level) *zap.SugaredLogger {
	atom := zap.NewAtomicLevelAt(level)
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.Sampling = nil
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.Level = atom

	zl, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("zap logger creation: %s", err))
	}

	return zl.Sugar()
}

func printVersion() {
	fmt.Printf("%s, version %s (revision: %s)\n\tbuild date: %s\n\tgo version: %s\n",
		binaryName,
		version,
		commit,
		date,
		runtime.Version(),
	)
}
