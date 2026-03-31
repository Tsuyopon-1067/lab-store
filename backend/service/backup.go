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

// StartIntervalBackupScheduler starts a goroutine that runs a backup every
// backup_interval_minutes as configured in the settings table. The interval
// is re-read from the database at the start of each sleep cycle so that
// changes to settings take effect on the next cycle (not mid-sleep).
func StartIntervalBackupScheduler(db *sql.DB, cfg *config.Config) {
	settingRepo := repository.NewSettingRepository(db)
	go func() {
		for {
			// Read current settings; fall back to config defaults if DB read fails
			intervalMinutes := cfg.Backup.Interval
			backupPath := cfg.Backup.Path

			setting, err := settingRepo.Get()
			if err != nil {
				log.Printf("backup scheduler: failed to read settings, using defaults: %v", err)
			} else {
				intervalMinutes = setting.BackupIntervalMinutes
				backupPath = setting.BackupPath
			}

			// Enforce safe minimum (10 minutes)
			const minIntervalMinutes = 10
			if intervalMinutes < minIntervalMinutes {
				log.Printf("backup scheduler: interval %d too small, using %d minutes",
					intervalMinutes, minIntervalMinutes)
				intervalMinutes = minIntervalMinutes
			}

			sleepDuration := time.Duration(intervalMinutes) * time.Minute
			log.Printf("backup scheduler: next run in %d minutes", intervalMinutes)

			time.Sleep(sleepDuration)

			// Create service with current backup path for this cycle
			svc := &BackupService{
				db:         db,
				backupPath: backupPath,
				dbPath:     cfg.Database.Path,
				repo:       repository.NewBackupRepository(db),
			}
			backup, err := svc.CreateBackup()
			if err != nil {
				log.Printf("auto backup failed: %v", err)
			} else {
				log.Printf("auto backup created: %s", backup.Filename)
			}
		}
	}()
}
