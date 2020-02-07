package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kubernetes-misc/chpa/model"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	asv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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
	logrus.Debugln("== getting namespaces ==")
	ls, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	result := make([]string, len(ls.Items))
	for i, n := range ls.Items {
		result[i] = n.Name
	}
	return result, nil
}

func GetAllCRD(namespace string, crd schema.GroupVersionResource) (result []model.CronHPAV1, err error) {
	logrus.Debugln("== getting CRDs ==")
	crdClient := dynClient.Resource(crd)
	ls, err := crdClient.Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorln(fmt.Errorf("could not list %s", err))
		return
	}

	result = make([]model.CronHPAV1, len(ls.Items))
	for i, entry := range ls.Items {
		b, err := entry.MarshalJSON()
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		cs := model.CronHPAV1{}
		err = json.Unmarshal(b, &cs)
		if err != nil {
			logrus.Errorln(err)
		}
		result[i] = cs
	}
	return
}

func GetDeployment(ns, name string) (deployment *v1.Deployment, err error) {
	logrus.Debugln("== getting deployment ==")
	deployment, err = clientset.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorln(err)
		return
	}
	return
}

func UpdateDeployment(deployment *v1.Deployment) (err error) {
	logrus.Debugln("== update deployment ==")
	_, err = clientset.AppsV1().Deployments(deployment.Namespace).Update(deployment)
	return
}

func GetHPA(ns, name string) (hpa *asv1.HorizontalPodAutoscaler, err error) {
	logrus.Debugln("== getting HPA ==")
	hpa, err = clientset.AutoscalingV1().HorizontalPodAutoscalers(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorln(err)
		return
	}
	return
}

func UpdateHPA(ns string, hpa *asv1.HorizontalPodAutoscaler) (err error) {
	logrus.Debugln("== updating HPA ==")
	_, err = clientset.AutoscalingV1().HorizontalPodAutoscalers(ns).Update(hpa)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	return
}
