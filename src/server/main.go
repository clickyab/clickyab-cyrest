package main

import (
	"common/config"
	"common/controllers/base"
	"common/models"
	"common/redis"
	"common/version"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize()

	if config.Config.DevelMode {
		// In development mode I need colors :) candy mode is GREAT!
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableColors: false})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: false, DisableColors: true})
		logrus.SetLevel(logrus.WarnLevel)
	}

	numcpu := config.Config.MaxCPUAvailable
	if numcpu < 1 || numcpu > runtime.NumCPU() {
		numcpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(numcpu)

	// Set global timezone
	if l, err := time.LoadLocation(config.Config.TimeZone); err == nil {
		time.Local = l
	}

	aredis.Initialize()
	models.Initialize()
	ver := version.GetVersion()

	logrus.WithFields(
		logrus.Fields{
			"Commit hash":       ver.Hash,
			"Commit short hash": ver.Short,
			"Commit date":       ver.Date.Format(time.RFC3339),
			"Build date":        ver.BuildDate.Format(time.RFC3339),
		},
	).Infof("Application started")

	base.Initialize(config.Config.MountPoint).Start(config.Config.Port)
}
