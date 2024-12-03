package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var (
	usedResourceNames = make(map[string]bool)
)

func GetUniqueName(prefix string) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = letters[rand.Intn(len(letters))]
	}
	name := fmt.Sprintf("%s%s", prefix, string(buf))
	if usedResourceNames[name] {
		return GetUniqueName(prefix)
	}
	usedResourceNames[name] = true
	return name
}

func doMain() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}

	client, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}

	inputFilePath := "job1/input.yaml"
	outputFilePath := "job1/output.yaml"
	name := GetUniqueName("job-")
	namespace := "default"

	inputYaml, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}
	var job batchv1.Job
	if err := yaml.Unmarshal(inputYaml, &job); err != nil {
		return err
	}
	job.SetName(name)
	job.SetNamespace(namespace)
	if err := client.Create(context.Background(), &job); err != nil {
		return err
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for i := 0; i < 10; i++ {
		var job batchv1.Job
		if err := client.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, &job); err != nil {
			return err
		}

		b, err := json.Marshal(job.Status)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(outputFile, "%s\n", b); err != nil {
			return err
		}
		fmt.Printf("%s\n", b)

		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {
	if err := doMain(); err != nil {
		log.Fatal(err)
	}
}
