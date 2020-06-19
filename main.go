package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/postfinance/kuota-calc/internal/calc"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
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

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
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
			log.Fatal(err)
		}

		fmt.Printf("CPU: %s\nMemory: %s\nOverhead: %f%%\n",
			usage.CPU,
			usage.Memory,
			usage.Overhead,
		)
	default:
		log.Fatalf("%s is not supported\n", gvk)
	}
}
