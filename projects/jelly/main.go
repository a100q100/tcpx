package main

import (
	"github.com/gin-gonic/gin"
	"tcpx/projects/jelly/control/config/configRouter"
)

func main() {
	go Http()
	go Tcpx()
	select {}
}

func Http() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pone")
	})
	configRouter.ConfigRouter(r)
	// userRouter.UserRouter(r)
}
func Tcpx() {
}
