package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"purchase-system/config"
	"purchase-system/model"
	"purchase-system/repository"
)

type BackupService struct {
	db         *sql.DB
	backupPath string
	dbPath     string
	repo       *repository.BackupRepository
}

func NewBackupService(db *sql.DB, cfg *config.Config) *BackupService {
	return &BackupService{
		db:         db,
		backupPath: cfg.Backup.Path,
		dbPath:     cfg.Database.Path,
		repo:       repository.NewBackupRepository(db),
	}
}

func (s *BackupService) CreateBackup() (*model.Backup, error) {
	// Ensure backup directory exists
	if err := os.MkdirAll(s.backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("backup-%s.db", timestamp)
	destPath := filepath.Join(s.backupPath, filename)

	// Execute VACUUM INTO to safely backup the database (works with WAL mode)
	vacuumSQL := fmt.Sprintf("VACUUM INTO '%s'", destPath)
	if _, err := s.db.Exec(vacuumSQL); err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Record backup in database
	backup, err := s.repo.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to record backup: %w", err)
	}

	return backup, nil
}

func (s *BackupService) ListBackups() ([]*model.Backup, error) {
	return s.repo.ListAll()
}

// StartDailyBackupScheduler starts a goroutine that runs a backup every day at midnight (local time).
func StartDailyBackupScheduler(db *sql.DB, cfg *config.Config) {
	svc := NewBackupService(db, cfg)
	go func() {
		for {
			// Calculate time until next midnight
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			waitDuration := time.Until(nextMidnight)

			// Log next execution time
			log.Printf("daily backup scheduler: next run at %s", nextMidnight.Format("2006-01-02 15:04:05"))

			// Wait until midnight
			time.Sleep(waitDuration)

			// Execute backup
			backup, err := svc.CreateBackup()
			if err != nil {
				log.Printf("auto backup failed: %v", err)
			} else {
				log.Printf("auto backup created: %s", backup.Filename)
			}
		}
	}()
}
