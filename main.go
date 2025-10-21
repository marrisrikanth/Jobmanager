package main

import (
	"jobmanager/api"
	"jobmanager/storage"
	"jobmanager/scheduler"
	"net/http"
	"os"
	"time"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	store, err := storage.NewSQLiteStore("jobmanager.db")
	if err != nil {
		log.Fatal("Failed to open Job store", err.Error)
	}

	s := scheduler.NewScheduler(store, log, 75.0)

	handler := &api.Handler{Store: store, Log: log, Scheduler: s,}

	//background purger
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			store.PurgeAndArchiveOldJobs(24 * time.Hour, log)
			<- ticker.C
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateJob(w, r)
		case http.MethodGet:
			handler.ListJobs(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/job", handler.GetJob)
	mux.HandleFunc("/job/cancel", handler.CancelJob)

	// Protect routes with API key middleware
	protected := api.APIKeyMiddleware(mux)
	//log.Info("protected:", protected)
	log.Info("Job Manager running on http://localhost:8080")
	http.ListenAndServe(":8080", protected)
	//http.ListenAndServe(":8443", "certs/server.crt", "certs/server.key", protected)
}
