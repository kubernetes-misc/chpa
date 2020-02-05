package model

type Job struct {
	F func()
}

func (j Job) Run() {
	j.F()
}
