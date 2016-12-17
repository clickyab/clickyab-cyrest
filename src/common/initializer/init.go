package initializer

import "github.com/Sirupsen/logrus"

// Initializer is the type to call early on system initialize call
type Initializer interface {
	Initialize()
}

// Finalizer is the type to call at finalize the object
type Finalizer interface {
	Finalize()
}

type group []interface{}

var (
	gr = make(group, 0)
)

func (g group) Finalize() {
	for i := range g {
		if in, ok := g[i].(Finalizer); ok {
			in.Finalize()
		}
	}
}

// Register a module in initializer
func Register(in interface{}) {

	switch in.(type) {
	case Initializer:
		gr = append(gr, in)
	case Finalizer:
		gr = append(gr, in)
	default:
		logrus.Panic("nor initializer nor finalizer")
	}
}

// Initialize all modules
func Initialize() Finalizer {
	for i := range gr {
		if in, ok := gr[i].(Initializer); ok {
			in.Initialize()
		}
	}

	return gr
}
