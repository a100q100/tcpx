package configRouter

import "github.com/gin-gonic/gin"

func ConfigRouter(r gin.IRoutes) {
	r.POST("/jelly/add-config/", configService.AddConfig)
}
