package misc

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/utils"

	"github.com/gin-gonic/gin"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *gin.Engine, mountPoint string) {

	groupMiddleware := gin.HandlersChain{}

	group := r.Group(mountPoint+"/misc", groupMiddleware...)

	// Route {/version GET Controller.getVersion misc []  Controller u  } with key 0
	m0 := gin.HandlersChain{}

	m0 = append(m0, u.getVersion)
	group.GET("/version", m0...)
	// End route {/version GET Controller.getVersion misc []  Controller u  } with key 0

	utils.DoInitialize(u)
}
