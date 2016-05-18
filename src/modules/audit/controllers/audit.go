package controllers

import (
	"common/controllers/base"
	"common/hub"

	"github.com/Sirupsen/logrus"
	"github.com/olebedev/emitter"
)

// Controller is the controller for the audit package
// @Route {
//		group = /audit
// }
type Controller struct {
}

// Initialize the controller
func (ctrl *Controller) Initialize() {
	go ctrl.auditLoop(hub.Subscribe("audit"))
}

func (ctrl *Controller) auditLoop(c <-chan emitter.Event) {
	for a := range c {
		logrus.Infof("Audit recieced %+v", a.Args)
	}
}

func init() {
	base.Register(&Controller{})
}
