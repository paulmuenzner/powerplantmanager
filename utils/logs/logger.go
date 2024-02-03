package log

import (
	"fmt"
	"github.com/paulmuenzner/powerplantmanager/config"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(logFileName string) {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logFolder := "log"
	err := os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating log directory:", err)
	}

	// Cleanup old log files
	cleanupOldLogFiles(logFolder)

	logFilePath := filepath.Join(logFolder, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func cleanupOldLogFiles(logFolder string) {
	files, err := filepath.Glob(filepath.Join(logFolder, "app_*.log"))
	if err != nil {
		log.Warn("Error globbing log files:", err)
		return
	}

	deleteLogsAfterDays := config.DeleteLogsAfterDays
	// Convert days to duration
	deleteLogsDuration := time.Duration(deleteLogsAfterDays) * 24 * time.Hour

	daysAgoToDeleteAfter := time.Now().Add(-1 * deleteLogsDuration)

	for _, file := range files {
		// Extract date from the filename using a regular expression
		re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})\.log$`)
		matches := re.FindStringSubmatch(file)
		if len(matches) != 2 {
			continue
		}

		dateString := matches[1]
		fileTime, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			log.Warn("Error parsing file date:", err)
			continue
		}

		// Check if the file is older than ten days
		if fileTime.Before(daysAgoToDeleteAfter) {
			err := os.Remove(file)
			if err != nil {
				log.Warn("Error deleting old log file:", err)
			} else {
				log.Infof("Deleted old log file: %s", file)
			}
		}
	}
}

func GetLogFileName() string {
	now := time.Now().Format("2006-01-02")
	return fmt.Sprintf("app_%s.log", now)
}

// GetLogger returns the configured logger
func GetLogger() *logrus.Logger {
	return log
}
