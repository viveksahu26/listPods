package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	log.Println(" KUBECONFIG flag is parsed: ", *kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods("monitoring").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s \n", pod.Name)
	}

	restClient := clientset.CoreV1().RESTClient()

	url := "/api/v1/namespaces/vivek-worspace/pods/vivek"
	fmt.Println("URL: ", url)

	response := restClient.Get().RequestURI(url).Do(context.TODO())

	// Extract the response body as a byte array
	body, err := response.Raw()
	if err != nil {
		panic(err.Error())
	}
	// Convert the byte array to a string
	responseString := string(body)
	fmt.Println("\n Pod Respons: ", responseString)
}

type Pod struct {
	Kind       string            `json:"kind"`
	ApiVersion string            `json:"apiVersion"`
	Metadata   map[string]string `json:"metadata"`
	// add more fields here if needed
}
