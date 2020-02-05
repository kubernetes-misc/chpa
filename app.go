package main

import (
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/client"
	"k8s.io/mouse/cron"
	"k8s.io/mouse/model"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	err := client.BuildClient()
	if err != nil {
		panic(err)
	}

	c, err := client.GetAllCRD("default", model.CronScaleV1CRDSchema)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	for cs := range c {
		cron.AddJobIfNotExists(cs)
	}

}
