package cron

import (
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/controller"
	"k8s.io/mouse/model"
)

var Jobs = make(map[string]cron.EntryID)

type CronScaleJob struct {
	cs model.CronScaleV1
}

func (j CronScaleJob) Run() {
	controller.ReconHub.Add(j.cs)
}

func AddJobIfNotExists(cs model.CronScaleV1) {
	_, found := Jobs[cs.GetID()]
	if found {
		return
	}

	entryID, err := cron.New(cron.WithSeconds()).AddJob(cs.Spec.CronSpec, Job{
		cs: cs,
	})
	if err != nil {
		logrus.Errorln(err)
		return
	}
	Jobs[cs.GetID()] = entryID
}
