package misc

import (
	"common/hub"
	"modules/misc/base"

	"github.com/Sirupsen/logrus"
	"github.com/olebedev/emitter"
)

// Controller is the misc controller
// @Route {
//		group = /misc
// }
type Controller struct {
	base.Controller
}

// Initialize is the initialization method for ths controller
// TODO : the trans module is somehow dangled in the system find its place!
func (u *Controller) Initialize() {
	//try.CatchHook(func(e error) error {
	//	return fmt.Errorf(trans.T(e.Error()))
	//})
	go u.auditLoop(hub.Subscribe("audit"))

}

func (u *Controller) auditLoop(c <-chan emitter.Event) {
	for a := range c {
		logrus.Infof("Audit recieced %+v", a.Args)
	}
}

func init() {
	base.Register(&Controller{})
}
