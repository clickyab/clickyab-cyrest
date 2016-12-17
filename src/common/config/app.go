package config

import (
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
)

// InitApplication is the helper function to initialize logger and some generic options
func InitApplication() {

	if Config.DevelMode {
		// In development mode I need colors :) candy mode is GREAT!
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableColors: false})
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: false, DisableColors: true})
		logrus.SetLevel(logrus.WarnLevel)
	}

	numcpu := Config.MaxCPUAvailable
	if numcpu < 1 || numcpu > runtime.NumCPU() {
		numcpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(numcpu)

	// Set global timezone
	if l, err := time.LoadLocation(Config.TimeZone); err == nil {
		time.Local = l
	}
}
