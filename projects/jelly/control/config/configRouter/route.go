package configRouter

import (
	"github.com/gin-gonic/gin"
	"tcpx/projects/jelly/control/config/configService"
)

func ConfigRouter(r gin.IRoutes) {
	r.POST("/jelly/add-config/", configService.HTTPAddConfig)
}
