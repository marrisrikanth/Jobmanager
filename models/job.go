package models

import "time"

type JobStatus string

const (
	Pending   JobStatus = "PENDING"
	Running   JobStatus = "RUNNING"
	Completed JobStatus = "COMPLETED"
	Failed    JobStatus = "FAILED"
	Cancelled JobStatus = "CANCELLED"
	Queued 	  JobStatus = "QUEUED"
)

type Job struct {
	ID        string    `json:"id" bson:"id"`
	Command   string    `json:"command" bson:"command"`
	Status    JobStatus `json:"status" bson:"status"`
	Output    string    `json:"output" bson:"output"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
