package main

import (
	"context"
	"flag"
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func getConfig() *rest.Config {
	kubeConfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	if *kubeConfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
		if err != nil {
			panic(err.Error())
		}
		return config
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

func main() {
	config := getConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc := obj.(*v1.Service)
			if svc.Spec.Type == "LoadBalancer" {
				ip := svc.Spec.LoadBalancerIP
				fmt.Printf("Setting %s to %s/%s\n", ip, svc.Namespace, svc.Name)
				svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{{IP: ip}}
				clientset.CoreV1().Services(svc.Namespace).UpdateStatus(context.TODO(), svc, metav1.UpdateOptions{})
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			svc := new.(*v1.Service)
			if svc.Spec.Type == "LoadBalancer" {
				ip := svc.Spec.LoadBalancerIP
				fmt.Printf("Setting %s to %s/%s\n", ip, svc.Namespace, svc.Name)
				svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{{IP: ip}}
				clientset.CoreV1().Services(svc.Namespace).UpdateStatus(context.TODO(), svc, metav1.UpdateOptions{})
			}
		},
		DeleteFunc: func(obj interface{}) {
		},
	}
	watcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "services", v1.NamespaceAll, fields.Everything())
	_, informer := cache.NewIndexerInformer(watcher, &v1.Service{}, 0, handlers, cache.Indexers{})

	informer.Run(nil)
}
