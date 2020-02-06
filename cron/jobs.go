package cron

import (
	"fmt"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/controller"
	"k8s.io/mouse/model"
)

var Jobs = make(map[string]*cron.Cron)

type CronScaleJob struct {
	cs model.CronScaleV1
}

func (j CronScaleJob) Run() {
	controller.ReconHub.Add(j.cs)
}

func MatchJobs(all []model.CronScaleV1) {
	logrus.Println(fmt.Sprintf("> Matching %v jobs", len(all)))
	//Find invalid jobs
	toRemove := make([]string, 0)
	for csID, _ := range Jobs {
		if IDExists(csID, all) {
			logrus.Println("... no problem with", csID)
			continue
		}
		logrus.Println("... should remove", csID)
		toRemove = append(toRemove, csID)
	}
	//Remove invalid jobs
	for _, tr := range toRemove {
		logrus.Println(fmt.Sprintf("> Removed job %s", tr))
		Jobs[tr].Stop()
		delete(Jobs, tr)
	}

	//Create if not exists
	for _, cs := range all {
		_, found := Jobs[cs.GetID()]
		if found {
			logrus.Println("...", cs.GetID(), "not creating as it already exists")
			continue
		}
		logrus.Println("...", cs.GetID(), "does not exist yet")
		c := cron.New()
		_, err := c.AddJob(cs.Spec.CronSpec, CronScaleJob{
			cs: cs,
		})
		c.Start()
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		logrus.Println(fmt.Sprintf("> Creating job %s", cs.GetID()))
		Jobs[cs.GetID()] = c
	}

}

func IDExists(id string, list []model.CronScaleV1) bool {
	for _, l := range list {
		if l.GetID() == id {
			return true
		}
	}
	return false
}
