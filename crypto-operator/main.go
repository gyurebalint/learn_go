package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for Namespace events... (Create one in Docker Desktop to see!)")
	watcher, err := clientset.CoreV1().Namespaces().Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for event := range watcher.ResultChan() {
		ns, ok := event.Object.(*v1.Namespace)
		if !ok {
			fmt.Println("unexpected type")
			continue
		}

		fmt.Printf("event deteted! Type %s | Namepspace: %s\n", event.Type, ns.Name)
	}
}
