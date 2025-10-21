package storage

import (
	"jobmanager/models"
	"sync"
)

type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*models.Job
}

func NewJobStore() *JobStore {
	return &JobStore{jobs: make(map[string]*models.Job)}
}

func (s *JobStore) Save(job *models.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

func (s *JobStore) Get(id string) (*models.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	return job, ok
}

func (s *JobStore) List() []*models.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*models.Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}
	return jobs
}
