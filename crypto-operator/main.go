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
			deployAggregator(clientset, deploy.Namespace)
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
				deployAggregator(clientset, ns.Name)
			}
		case "DELETED":
			fmt.Printf("namespace %s was deleted. Cleaning up local", ns.Name)
		case "MODIFIED":
			fmt.Printf("Namespace %s was updated.\n", ns.Name)
		}

		fmt.Printf("event deteted! Type %s | Namepspace: %s\n", event.Type, ns.Name)
	}
}

func deployPostgres(clientset *kubernetes.Clientset, nsName string) error {
	//check if postgres deployment exists
	_, err := clientset.AppsV1().Deployments(nsName).Get(context.TODO(), "postgres-db", metav1.GetOptions{})
	if err == nil {
		fmt.Println("🐘 Postgres Deployment already exists. Skipping creation.")
	} else {
		fmt.Printf("🐘 Spinning up Postgres in %s...\n", nsName)

		dep := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: "postgres-db"},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "postgres"}},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "postgres"}},
					Spec: v1.PodSpec{
						Containers: []v1.Container{{
							Name:  "postgres",
							Image: "postgres:15-alpine",
							Env: []v1.EnvVar{
								{Name: "POSTGRES_USER", Value: "admin"},
								{Name: "POSTGRES_DB", Value: "crypto-aggregator"},
								{Name: "POSTGRES_PASSWORD", Value: "admin"}, // In prod, use a Secret!
							},
							Ports: []v1.ContainerPort{{ContainerPort: 5432}},
						}},
					},
				},
			},
		}
		_, err := clientset.AppsV1().Deployments(nsName).Create(context.TODO(), dep, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	_, err = clientset.CoreV1().Services(nsName).Get(context.TODO(), "postgres-db", metav1.GetOptions{})
	if err == nil {
		fmt.Println("networks Postgres Service already exists. Skipping creation.")
		return nil
	}
	
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "postgres-db"},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": "postgres"},
			Ports:    []v1.ServicePort{{Port: 5432}},
		},
	}
	_, err = clientset.CoreV1().Services(nsName).Create(context.TODO(), svc, metav1.CreateOptions{})
	return err
}

func deployAggregator(clientset *kubernetes.Clientset, nsName string) {
	err := deployPostgres(clientset, nsName)
	if err != nil {
		fmt.Printf("❌ Failed to deploy Postgres: %v\n", err)
		return
	}

	fmt.Printf("🚀 Deploying Aggregator v2 to %s...\n", nsName)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "crypto-aggregator"},
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
							Image: "crypto-aggregator:v2", // Using your new image
							Env: []v1.EnvVar{
								{Name: "DB_HOST", Value: "postgres-db"}, // Matches the PG Service name
								{Name: "DB_USER", Value: "admin"},
								{Name: "DB_NAME", Value: "crypto-aggregator"},
								{Name: "DB_PASSWORD", Value: "admin"},
							},
						},
					},
				},
			},
		},
	}

	_, err = clientset.AppsV1().Deployments(nsName).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("❌ Failed to deploy Aggregator: %v\n", err)
	}
}

func int32Ptr(i int32) *int32 { return &i }
