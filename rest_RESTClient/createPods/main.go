package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string

	// check if machine has home directory.
	if home := homedir.HomeDir(); home != "" {
		// read kubeconfig flag. if not provided use config file $HOME/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	log.Println("KUBECONFIG flag is parsed: ", *kubeconfig)

	// build configuration from the config file.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// create kubernetes clientset. this clientset can be used to create,delete,patch,list etc for the kubernetes resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// build the pod defination we want to deploy
	// generate a pod definition that we want to deploy.
	pod := getPodObject()

	// now create the pod in kubernetes cluster using the clientset
	// create the pod in kubernetes cluster using the clientset.
	pod, err = clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s not found in %s namespace\n", pod.Name, pod.Namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error creating %s %v\n", statusError.ErrStatus.Message, pod.Name)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found %s pod in %s namespace\n", pod.Name, pod.Namespace)
	}
	fmt.Printf("%s Created successfully... ", pod.Name)
	fmt.Println("Pod created successfully...: ", pod)
}

func getPodObject() *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-test-pod",
			Namespace: "vivek-worspace",
			Labels: map[string]string{
				"app": "demo",
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "busybox",
					Image:           "busybox",
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
				},
			},
		},
	}
}
