package model

import "time"

type Job interface {
	Execute()
}

type Worker interface {
	Run()
	Stop()
	JobChannel() chan<- Job
}

type Scheduler interface {
	Name() string
	JobType() string
	Enabled(cfg *Config) bool
	NextScheduleTime(cfg *Config, now time.Time, pendingJobs bool, lastSuccessfulJob Job) *time.Time
	ScheduleJob(cfg *Config, pendingJobs bool, lastSuccessfulJob Job) (Job, *AppError)
}
