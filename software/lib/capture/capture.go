package capture

import (
	"fmt"
	"io"
	"main/lib/config"
	"main/lib/shared"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

type CaptureManager struct {
	Logger *zap.Logger
	Config *config.Config
}

func (c CaptureManager) Capture() {
	currentTime := time.Now().UTC()
	timestamp := currentTime.Format(shared.DATE_FORMAT)
	c.Logger.Info("Capture started", zap.String("timestamp", timestamp))

	for _, host := range c.Config.CaptureConfig.CaptureHosts {
		go func(host config.Server) {
			hLog := c.Logger.With(zap.Any("host", host))
			filename := filepath.Join(c.Config.CaptureConfig.CaptureDirectory, host.Server_Name, fmt.Sprintf("%s.jpg", timestamp))

			// Download the image
			if err := downloadImage(host.Server_Url, filename); err != nil {
				hLog.Error("Failed to download image", zap.Error(err))
			} else {
				hLog.Info("Image saved", zap.String("filename", filename))
			}
		}(host)
	}
}

func downloadImage(url, filepath string) error {
	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
