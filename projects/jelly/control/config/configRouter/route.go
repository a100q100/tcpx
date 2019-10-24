package configRouter

import (
	"github.com/gin-gonic/gin"
	"tcpx/projects/jelly/control/config/configService"
)

func ConfigRouter(r gin.IRoutes) {
	r.POST("/jelly/config/", configService.HTTPAddConfig)
	r.DELETE("/jelly/config/:id/", configService.HTTPDeleteConfig)
	r.GET("/jelly/config/", configService.HTTPListConfig)
	r.GET("/jelly/config/:id/", configService.HTTPGetConfig)
	r.PATCH("/jelly/config/:id", configService.HTTPUpdateConfig)
}
