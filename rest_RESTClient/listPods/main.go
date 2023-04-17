package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getPods() {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++INSIDE Get PODS +++++++++++++++++++++++++++++++++++")
	defer fmt.Println("--------------------------------EXIT get PODS ----------------------------------")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Print("Failed to instantiate k8s client: ", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	restClient := clientset.CoreV1().RESTClient()
	url1 := "/api/v1/namespaces/default/pods"
	fmt.Println("URL1: ", url1)
	response1 := restClient.Get().RequestURI(url1).Do(context.TODO())
	// Extract the response body as a byte array
	body1, err := response1.Raw()
	if err != nil {
		panic(err.Error())
	}
	// Convert the byte array to a string
	responseString1 := string(body1)
	fmt.Println("\n PODS IN DEFAULT: ", responseString1)

	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods("obmondo").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("FAILED TO LIST PODS INSIDE OBMONDO")
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s \n", pod.Name)
	}

	fmt.Println("*******************************************************************************************************************")
	// url := "/apis/argoproj.io/v1alpha1/namespaces/argocd/applications"
	// fmt.Println("URL: ", url)

	// response := restClient.Get().RequestURI(url).Do(context.TODO())
	fmt.Println("ARGOCD v1alpha1 USECASE")
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating dynamic client:", err)
	}
	// Define the GVK (Group, Version, Kind) of the Argo CD Application resource
	applicationGVK := schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "Application",
	}
	fmt.Println("applicationGVK: ", applicationGVK)
	applications, err := dynamicClient.Resource(applicationGVK).Namespace("argocd").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Print the names of all applications
	for _, app := range applications.Items {
		metaObj, err := meta.Accessor(&app)
		if err != nil {
			panic(err)
		}
		fmt.Println("NAME:", metaObj.GetName())
	}

	// // Extract the response body as a byte array
	// body, err := response.Raw()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// // Convert the byte array to a string
	// responseString := string(body)
	// fmt.Println("\n APPS IN ARGOCD: ", responseString)
	fmt.Println("*******************************************************************************************************************")

	argocdClientset, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Argo CD clientset: %v", err)
	}
	// List all Argo CD applications
	argoCDAppsList, err := argocdClientset.ArgoprojV1alpha1().Applications("argocd").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list Argo CD applications: %v", err)
	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	for i, app := range argoCDAppsList.Items {
		appName := app.GetName()
		applicationNamespace := app.Spec.Destination.Namespace
		fmt.Printf("%v ARGOCD APP NAME: %s, NAMESPACE: %s, LABELS: %s \n", i, appName, applicationNamespace, app.Labels)

	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}

func runner(ticker *time.Ticker, done <-chan bool) {
	for {
		select {
		case <-done:
			return
		// this will get triggered after a particular duration of time.
		case <-ticker.C:
			fmt.Println("Getting pods...")
			getPods()

		}
	}
}

func main() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// channel to mark completion
	done := make(chan bool)

	// await termination signals from OS on a channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go runner(ticker, done)

	// blocks the main goroutine infinitely until user terminates the process
	sig := <-shutdown
	done <- true
	log.Printf("received %s, terminating application", sig)
	close(shutdown)
	close(done)
	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()
	// log.Println(" KUBECONFIG flag is parsed: ", *kubeconfig)

	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	panic(err)
	// }
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Print("Failed to instantiate k8s client: ", err)
	// 	os.Exit(1)
	// }
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	// // get pods in all the namespaces by omitting namespace
	// // Or specify namespace to get pods in particular namespace
	// pods, err := clientset.CoreV1().Pods("monitoring").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// for _, pod := range pods.Items {
	// 	fmt.Printf("Pod Name: %s \n", pod.Name)
	// }

	// restClient := clientset.CoreV1().RESTClient()

	// url := "/apis/argoproj.io/v1alpha1/namespaces/argocd/applications"
	// fmt.Println("URL: ", url)

	// response := restClient.Get().RequestURI(url).Do(context.TODO())

	// // Extract the response body as a byte array
	// body, err := response.Raw()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// // Convert the byte array to a string
	// responseString := string(body)
	// fmt.Println("\n Pod Respons: ", responseString)
}

type Pod struct {
	Kind       string            `json:"kind"`
	ApiVersion string            `json:"apiVersion"`
	Metadata   map[string]string `json:"metadata"`
	// add more fields here if needed
}
