package delete

import (
	"fmt"
	"main/lib/config"
	"main/lib/shared"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

type DeleteManager struct {
	Log    *zap.Logger
	Config *config.Config
}

func (d DeleteManager) DeleteOldFiles() {
	currentTime := time.Now().UTC()
	d.Log.Info("Delete started", zap.Time("time", currentTime))
	dirPath := d.Config.CaptureConfig.CaptureDirectory
	delete_ms := d.Config.DeleteConfig.DeleteMS

	// Walk through the directory, but only one level deep
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			d.Log.Error("Failed to access path", zap.Error(err))
			return err
		}

		// Only look at files in subdirectories (skip the top-level directory)
		if info.IsDir() && path != dirPath {
			files, err := os.ReadDir(path)
			if err != nil {
				d.Log.Error("Failed to read subdirectory", zap.Error(err))
				return err
			}

			// Iterate through files in the subdirectory
			for _, file := range files {
				if !file.IsDir() {
					fileName := file.Name()
					d.Log.Debug("next file", zap.Any("name", path+"/"+fileName))
					timestampStr := strings.TrimSuffix(fileName, filepath.Ext(fileName))

					parsedTime, err := time.Parse(shared.DATE_FORMAT, timestampStr)
					if err != nil {
						d.Log.Warn(fmt.Sprintf("Unparseable filename: %s", fileName))
						continue
					}
					d.Log.Debug("Parsed time", zap.Any("input", timestampStr), zap.Any("output", parsedTime))

					// Check if the file is older than the threshold
					timeDifference := currentTime.Sub(parsedTime).Milliseconds()
					d.Log.Debug("time difference",
						zap.Any("current", currentTime),
						zap.Any("currentFormatted", currentTime.Format(shared.DATE_FORMAT)),
						zap.Any("parsed", parsedTime),
						zap.Any("diff", timeDifference),
						zap.Any("deletems", delete_ms))
					if timeDifference > delete_ms {
						d.Log.Debug("deleting file")
						// Delete the file
						fullFilePath := filepath.Join(path, fileName)
						if err := os.Remove(fullFilePath); err != nil {
							d.Log.Error("Failed to delete file", zap.Error(err))
						} else {
							d.Log.Info(fmt.Sprintf("Deleted file: %s", fullFilePath))
						}
					} else {
						d.Log.Debug("Keeping file")
					}
				}
			}
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		d.Log.Error("Error walking the directory", zap.Error(err))
	}
}
