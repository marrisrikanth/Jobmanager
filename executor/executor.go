package executor

import (
	"bytes"
	"context"
	"jobmanager/models"
	"jobmanager/storage"
	"os/exec"
	"time"
	"sync"
)

var (
	activeJobs = make(map[string]context.CancelFunc)
	activeJobsMu sync.Mutex
)

func RunJob(job *models.Job, store *storage.SQLiteStore) {
	ctx, cancel := context.WithCancel(context.Background())

	activeJobsMu.Lock()
	activeJobs[job.ID] = cancel
	activeJobsMu.Unlock()

	job.Status = models.Running
	job.UpdatedAt = time.Now()
	store.Save(job)

	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, "bash", "-c", job.Command)
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	activeJobsMu.Lock()
	delete(activeJobs, job.ID)
	activeJobsMu.Unlock()

	if ctx.Err() == context.Canceled {
		job.Status = models.Cancelled
		job.Output = "Job cancelled"
	} else if err != nil {
		job.Status = models.Failed
		job.Output = out.String() + "\nError:" + err.Error()
	} else {
		job.Status = models.Completed
		job.Output = out.String()
	}
	job.UpdatedAt = time.Now()
	store.Save(job)
}

func CancelJob(id string, store *storage.SQLiteStore) error {
	
	activeJobsMu.Lock()
	cancel, exists := activeJobs[id]
	activeJobsMu.Unlock()

	if !exists {
		return nil // job doesn't exist
	}
	cancel()

	job, err := store.Get(id)
	if err == nil {
		job.Status = models.Cancelled
		job.Output = "Cancelled by user"
		job.UpdatedAt = time.Now()
		store.Save(job)
	}

	return nil
}

