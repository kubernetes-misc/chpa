package main

import (
	cronV3 "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/client"
	"k8s.io/mouse/cron"
	"k8s.io/mouse/model"
	"os"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
)

func main() {
	err := client.BuildClient()
	if err != nil {
		panic(err)
	}
	cronSpec := os.Getenv("cronSpec")
	cronV3.New().AddJob(cronSpec, model.Job{
		F: updateCronScales,
	})

}

func updateCronScales() {
	c, err := client.GetAllCRD("default", model.CronScaleV1CRDSchema)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	for cs := range c {
		cron.AddJobIfNotExists(cs)
	}

}
