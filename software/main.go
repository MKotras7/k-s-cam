package main

import (
	"main/lib/capture"
	"main/lib/config"
	"main/lib/delete"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	logger       *zap.Logger
	systemConfig *config.Config
	ticker       *time.Ticker // Ticker can be controlled to adjust the capture rate
)

func main() {
	var err error

	// Start logger
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.SetLevel(zap.DebugLevel)
	logger, err = zapConfig.Build()
	if err != nil {
		panic("Failed to initialize logger")
	}
	defer logger.Sync() // flushes buffer, if any
	logger.Info("Logger started")

	// Load config
	loadedConfig, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config file", zap.Error(err))
	}
	systemConfig = loadedConfig

	// Create necessary directories
	err = createDirectories()
	if err != nil {
		logger.Fatal("Failed to create directories", zap.Error(err))
	}
	logger.Info("Confirmed directories")

	startTimer()
	// Wait forever while timer runs
	for {
		time.Sleep(time.Minute)
	}
}

func createDirectories() error {
	var err error
	err = os.MkdirAll(systemConfig.CaptureConfig.CaptureDirectory, os.ModePerm)
	if err != nil {
		logger.Error("Failed to create capture directory", zap.Error(err))
		return err
	}
	for _, host := range systemConfig.CaptureConfig.CaptureHosts {
		err = os.MkdirAll(systemConfig.CaptureConfig.CaptureDirectory+"/"+host.Server_Name, os.ModePerm)
		if err != nil {
			logger.Error("Failed to create capture directory", zap.Error(err))
			return err
		}
	}
	return nil
}

func startTimer() {
	captureManager := capture.CaptureManager{
		Logger: logger,
		Config: systemConfig,
	}
	captureManager.Capture()
	go func() {
		ticker = time.NewTicker(time.Millisecond * time.Duration(systemConfig.CaptureConfig.IntervalMS))
		for range ticker.C {
			captureManager.Capture()
		}
	}()

	deleteManager := delete.DeleteManager{
		Log:    logger,
		Config: systemConfig,
	}
	deleteManager.DeleteOldFiles()
	go func() {
		ticker = time.NewTicker(time.Millisecond * time.Duration(systemConfig.DeleteConfig.IntervalMS))
		for range ticker.C {
			deleteManager.DeleteOldFiles()
		}
	}()
}
