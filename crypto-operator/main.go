package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		watchNameSpaces(clientset)
	}()

	go func() {
		defer wg.Done()
		watchDeployments(clientset)
	}()

	wg.Wait()
}

func watchDeployments(clientset *kubernetes.Clientset) {
	fmt.Println("🛡️  Deployment Guard is active...")
	watcher, _ := clientset.AppsV1().Deployments("").Watch(context.Background(), metav1.ListOptions{})

	for event := range watcher.ResultChan() {
		deploy, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			continue
		}

		if event.Type == "DELETED" && deploy.Name == "crypto-aggregator" {
			ns, err := clientset.CoreV1().Namespaces().Get(context.Background(), deploy.Namespace, metav1.GetOptions{})
			if err != nil || ns.Status.Phase == v1.NamespaceTerminating {
				fmt.Println("namespace doesn't exist anymore, we don't deploy crypto-aggregator")
				continue
			}

			fmt.Printf("⚠️  SABOTAGE DETECTED! Re-deploying %s to %s...\n", deploy.Name, deploy.Namespace)
			err = deployAggregator(clientset, deploy.Namespace)
			if err != nil {
				fmt.Printf("Failed to deploy: %v\n", err)
			}
		}
	}
}

func watchNameSpaces(clientset *kubernetes.Clientset) {
	fmt.Println("Waiting for Namespace events... (Create one in Docker Desktop to see!)")
	watcher, err := clientset.CoreV1().Namespaces().Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for event := range watcher.ResultChan() {
		ns, ok := event.Object.(*v1.Namespace)
		if !ok {
			continue
		}
		switch event.Type {
		case "ADDED":
			if strings.HasPrefix(ns.Name, "crypto-") {
				fmt.Printf("new namespace detected: %s. Deploying...\n", ns.Name)
				err := deployAggregator(clientset, ns.Name)
				if err != nil {
					fmt.Printf("Failed to deploy: %v\n", err)
				}
			}
		case "DELETED":
			fmt.Printf("namespace %s was deleted. Cleaning up local", ns.Name)
		case "MODIFIED":
			fmt.Printf("Namespace %s was updated.\n", ns.Name)
		}

		fmt.Printf("event deteted! Type %s | Namepspace: %s\n", event.Type, ns.Name)
	}
}

func deployAggregator(clientset *kubernetes.Clientset, name string) error {
	_, err := clientset.AppsV1().Deployments(name).Get(context.Background(), "crypto-aggregator", metav1.GetOptions{})
	if err == nil {
		return nil
	}

	fmt.Printf("🚀 Deploying Aggregator to %s...\n", name)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crypto-aggregator",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "aggregator"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "aggregator"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "aggregator",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
	_, err = clientset.AppsV1().Deployments(name).Create(context.Background(), deployment, metav1.CreateOptions{})
	return err
}

func int32Ptr(i int32) *int32 { return &i }
