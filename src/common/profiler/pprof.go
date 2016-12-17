package profiler

import (
	"common/config"
	"common/initializer"

	"common/utils"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/profile"
)

type profilerInitializer struct {
	profiler interface {
		Stop()
	}
}

func (pi *profilerInitializer) Initialize() {
	switch config.Config.Profile {
	case "cpu":
		pi.profiler = profile.Start(
			profile.CPUProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(config.Config.ProfileRoot, <-utils.ID)),
		)
	case "mem":
		pi.profiler = profile.Start(
			profile.MemProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(config.Config.ProfileRoot, <-utils.ID)),
		)
	case "trace":
		pi.profiler = profile.Start(
			profile.TraceProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(config.Config.ProfileRoot, <-utils.ID)),
		)
	case "block":
		pi.profiler = profile.Start(
			profile.BlockProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(config.Config.ProfileRoot, <-utils.ID)),
		)
	default:
		logrus.Debug("Profiler disabled")
	}
}

func (pi *profilerInitializer) Finalize() {
	if pi.profiler != nil {
		pi.profiler.Stop()
		logrus.Debug("Profiler done")
	}
}

func init() {
	initializer.Register(&profilerInitializer{})
}
