package scheduler

import (
	"time"
	"jobmanager/models"
	"jobmanager/storage"
	"jobmanager/executor"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	Store *storage.SQLiteStore
	Log	*logrus.Logger
	JobQueue chan *models.Job
	Threshold float64
}

func NewScheduler(store *storage.SQLiteStore, log *logrus.Logger, threshold float64) *Scheduler {
	s := &Scheduler{
		Store: store,
		Log: log,
		JobQueue: make(chan *models.Job, 10),
		Threshold: threshold,
	}

	go s.monitorAndRun()
	return s
}

func (s *Scheduler) Submit(job *models.Job) {

	go func() {
		load,_ := getCPULoad()
		if load > s.Threshold {
			job.Status = models.Queued
			job.Output = "System load high, job queued"
			s.Store.Save(job)
			s.JobQueue <- job
			s.Log.WithFields(logrus.Fields{
				"job_id":job.ID,
				"load": load,
			}).Info("Job queued due to high CPU load")
		} else {
			s.Log.WithFields(logrus.Fields{
				"job_id": job.ID,
				"load": load,
			}).Info("Starting job immediately")
			go executor.RunJob(job, s.Store)
		}
	}()
}

func (s *Scheduler) monitorAndRun() {
	for {
		time.Sleep(5 * time.Second)

		load,_ := getCPULoad()
		if load < s.Threshold && len(s.JobQueue) > 0 {
			s.Log.WithField("load", load).Info("CPU load low - execute")

			select {
			case job := <- s.JobQueue:
				job.Status = models.Pending
				job.Output = "Dequeued and starting job"
				job.UpdatedAt = time.Now()
				s.Store.Save(job)
				go executor.RunJob(job, s.Store)
			default:
				//queue empty
			}
		}
	}
}

func getCPULoad() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	if len(percent) > 0 {
		return percent[0], nil
	}
	return 0, nil
}


