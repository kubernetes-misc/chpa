package cron

import (
	"fmt"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/controller"
	"k8s.io/mouse/model"
	"sync"
)

var Jobs = make(map[string]*Job)

type Job struct {
	sync.Mutex
	CronScale model.CronScaleV1
	Cron      *cron.Cron
}

func (j Job) Run() {
	controller.ReconHub.Add(j.CronScale)
}

func (j *Job) UpdateCronScale(cs model.CronScaleV1) {
	j.Lock()
	j.CronScale = cs
	j.Unlock()
}

func MatchJobs(all []model.CronScaleV1) {
	logrus.Debugln(fmt.Sprintf("> Matching %v jobs", len(all)))
	//Find invalid jobs
	toRemove := make([]string, 0)
	for csID, _ := range Jobs {
		if IDExists(csID, all) {
			logrus.Debugln("... no problem with", csID)
			continue
		}
		logrus.Debugln("... should remove", csID)
		toRemove = append(toRemove, csID)
	}
	//Remove invalid jobs
	for _, tr := range toRemove {
		logrus.Infoln(fmt.Sprintf("> Removed job %s", tr))
		Jobs[tr].Cron.Stop()
		delete(Jobs, tr)
	}

	//Create if not exists
	for _, cs := range all {
		foundCS, found := Jobs[cs.GetID()]
		if found {
			logrus.Debugln("...", cs.GetID(), "updating as already exists")
			foundCS.UpdateCronScale(cs)
			continue
		}
		logrus.Debugln("...", cs.GetID(), "should be created")
		c := cron.New()
		j := &Job{
			CronScale: cs,
			Cron:      c,
		}
		_, err := c.AddJob(cs.Spec.CronSpec, j)
		c.Start()
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		logrus.Infoln(fmt.Sprintf("> Creating job %s as:", cs.GetID()))
		logrus.Infoln(cs.PrettyString())
		Jobs[cs.GetID()] = j
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
