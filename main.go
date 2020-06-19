package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/postfinance/kuota-calc/internal/calc"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/deprecated/scheme"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if (fi.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println(`kuota-calc calculates resource quotas based on your k8s yamls

Usage: cat deployment.yaml | kuota-calc`)

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
			log.Fatalf("decode: %s\n", err)
		}

		switch obj := object.(type) {
		case *appsv1.Deployment:
			usage, err := calc.Deployment(*obj)
			if err != nil {
				log.Printf("ERROR: %s\n", err)
				continue
			}

			summary = append(summary, usage)
		case *appsv1.StatefulSet:
			usage, err := calc.StatefulSet(*obj)
			if err != nil {
				log.Printf("ERROR: %s\n", err)
				continue
			}

			summary = append(summary, usage)
		default:
			log.Printf("ignoring %s\n", gvk)
			continue
		}
	}

	var (
		cpuUsage    resource.Quantity
		memoryUsage resource.Quantity
	)

	for _, u := range summary {
		cpuUsage.Add(*u.CPU)
		memoryUsage.Add(*u.Memory)
	}

	fmt.Printf("CPU: %s\nMemory: %s\n",
		cpuUsage.String(),
		memoryUsage.String(),
	)
}
