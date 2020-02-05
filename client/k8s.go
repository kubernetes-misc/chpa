package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/mouse/model"
	"path/filepath"
)

var clientset *kubernetes.Clientset
var dynClient dynamic.Interface

func BuildClient() (err error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	dynClient, err = dynamic.NewForConfig(config)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Error received creating client %v", err))
		return
	}
	return
}

func GetAllNS() ([]string, error) {
	ls, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	result := make([]string, len(ls.Items))
	for _, n := range ls.Items {
		result = append(result, n.Name)
	}
	return result, nil
}

func GetAllCRD(namespace string, crd schema.GroupVersionResource) (result chan model.CronScaleV1, err error) {
	crdClient := dynClient.Resource(crd)
	ls, err := crdClient.Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorln(fmt.Errorf("could not list %s", err))
		return
	}

	result = make(chan model.CronScaleV1)
	go func() {
		for _, entry := range ls.Items {
			b, err := entry.MarshalJSON()
			if err != nil {
				logrus.Errorln(err)
				continue
			}
			cs := model.CronScaleV1{}
			err = json.Unmarshal(b, &cs)
			if err != nil {
				logrus.Errorln(err)
			}
			result <- cs
			//fmt.Println(fmt.Sprintf("%s replicas: %v ==> %v @ CPU load of %v%% (cronscale/%s operating on %s/%s)", pad(cs.Spec.CronSpec, 12), cs.Spec.MinReplicas, cs.Spec.MaxReplicas, cs.Spec.TargetCPUUtilizationPercentage, cs.Metadata.Name, cs.Spec.ScaleTargetRef.Kind, cs.Spec.ScaleTargetRef.Name))
		}
		close(result)
	}()
	return
}

func GetDeployment(ns, name string) (deployment *v1.Deployment, err error) {
	deployment, err = clientset.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorln(err)
		return
	}
	return
}

func UpdateDeployment(ns string, deployment *v1.Deployment) (err error) {
	deployment, err = clientset.AppsV1().Deployments(ns).UpdateScale(deployment, &v1.Scale{})
	if err != nil {
		logrus.Errorln(err)
		return
	}
	return
}
