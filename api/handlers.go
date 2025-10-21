package api

import (
	"encoding/json"
	"jobmanager/executor"
	"jobmanager/models"
	"jobmanager/scheduler"
	"jobmanager/storage"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Store *storage.SQLiteStore
	Log	*logrus.Logger
	Scheduler *scheduler.Scheduler
}

// POST /jobs
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := &models.Job{
		ID:        uuid.NewString(),
		Command:   req.Command,
		Status:    models.Pending,
		Output: "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := h.Store.Save(job)
	if err != nil {
		h.Log.Info("Error:", err)
	}
	h.Log.WithFields(logrus.Fields{
		"job_id": job.ID,
		"command": job.Command,
	}).Info("Job created")
	//Run in background
	//go executor.RunJob(job, h.Store)
	h.Scheduler.Submit(job)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// GET /jobs
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, _ := h.Store.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// GET /jobs?id=<job-id>
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	job, ok := h.Store.Get(id)
	if ok != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *Handler) CancelJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	err := executor.CancelJob(id, h.Store)
	if err != nil {
		http.Error(w, "failed to cancel job", http.StatusInternalServerError)
		return
	}

	h.Log.WithField("job_id", id).Info("Job cancelled by user")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Job cancelled successfully",
	})
}

