package storage

import (
	"fmt"
	"database/sql"
	"jobmanager/models"
	_ "modernc.org/sqlite"
	"time"
	"github.com/sirupsen/logrus"
)

type SQLiteStore struct {
	DB *sql.DB
}

func NewSQLiteStore(filepath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	//Enable foreign key enforcement
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	migrations := []Migration{
	{
			Version: 1,
			UpSQL: `
		CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			command TEXT NOT NULL,
			status TEXT NOT NULL,
			output TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);
		`,
	},
	{
		Version: 2,
		UpSQL: `
		CREATE TABLE IF NOT EXISTS job_dependencies (
			job_id TEXT NOT NULL,
			depends_on TEXT NOT NULL,
			FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
			FOREIGN KEY (depends_on) REFERENCES jobs(id) ON DELETE CASCADE,
			PRIMARY KEY (job_id, depends_on)
		);
		`,

	},
	{
		Version: 3,
		UpSQL: `
		ALTER TABLE jobs ADD COLUMN cpu_load REAL DEFAULT 0;
		`,
	},
	{
		Version: 4,
		UpSQL: `
		CREATE TABLE IF NOT EXISTS archived_jobs (
		id TEXT PRIMARY KEY,
		command TEXT,
		status TEXT,
		output TEXT,
		created_at DATETIME,
		archived_at DATETIME DEFAULT CURRENT_TIMESTAMP);
		`,
	},
	}

	if err := applyMigrations(db, migrations); err != nil {
		return nil, fmt.Errorf("failed migrations: %w", err)
	}

	return &SQLiteStore{DB: db}, nil
}

func (s *SQLiteStore) Save(job *models.Job) error {

	fmt.Println(job)
	_, err := s.DB.Exec(`
	INSERT INTO jobs (id, command, status, output, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?) 
	ON CONFLICT(id) DO UPDATE SET
	command=excluded.command,
	status=excluded.status,
	output=excluded.output,
	created_at=excluded.created_at,
	updated_at=excluded.updated_at
	`,
	job.ID, job.Command, job.Status, job.Output, job.CreatedAt, job.UpdatedAt)
	return err
}

func (s *SQLiteStore) Get(id string) (*models.Job, error) {
	row := s.DB.QueryRow("SELECT id, command, status, output, created_at, updated_at FROM jobs WHERE id = ?", id)
	job := models.Job{}
	err := row.Scan(&job.ID, &job.Command, &job.Status, &job.Output, &job.CreatedAt, &job.UpdatedAt)
	return &job, err
}


func (s *SQLiteStore) List() ([]*models.Job, error) {
	rows, err := s.DB.Query("SELECT id, command, status, output, created_at, updated_at FROM jobs")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job

	for rows.Next() {
		job := &models.Job{}

		err := rows.Scan(&job.ID, &job.Command, &job.Status, &job.Output, &job.CreatedAt, &job.UpdatedAt)
		if err == nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}


func (s *SQLiteStore) ListQueuedJobs() ([]*models.Job, error) {
	rows, err := s.DB.Query("SELECT id, command, status, output, created_at, updated_at FROM jobs WHERE status = ?", models.Queued)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job

	for rows.Next() {
		job := &models.Job{}

		err := rows.Scan(&job.ID, &job.Command, &job.Status, &job.Output, &job.CreatedAt, &job.UpdatedAt)
		if err == nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (s *SQLiteStore) PurgeAndArchiveOldJobs(olderThan time.Duration, log *logrus.Logger) {
	/*
	cutoff := time.Now().Add(-olderThan)

	tx, err := s.DB.Begin()
	if err != nil {
		log.WithError(err).Error("Failed to begin transaction for purge")
		return
	}
	defer tx.Rollback()

	//copy old jobs into archived_jobs
	_, err = tx.Exec(`
	INSERT INTO archived_jobs (id, command, status, output, created_at)
	SELECT id, command, status, output, created_at FROM jobs WHERE created_at < ?;
	`, cutoff)
	if err != nil {
		log.WithError(err).Error("Failed to archive old jobs")
		tx.Rollback()
		return
	}

	result, err := s.DB.Exec(`DELETE FROM jobs WHERE created_at < ?`, cutoff)
	if err != nil {
		log.WithError(err).Error("Failed to delete old jobs after archiving")
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error("failed to commit purge transaction")
		return
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		log.WithFields(logrus.Fields{
			"archived_count": rows,
			"cutoff_time": cutoff,
		}).Info("Archived and Purged old jobs from SQLite store")
	}*/
}

