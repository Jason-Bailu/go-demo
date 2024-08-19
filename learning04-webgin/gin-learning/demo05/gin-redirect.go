package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// http redirect
	r.GET("/http/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})

	// route redirect
	r.GET("/route/redirect", func(c *gin.Context) {
		// 指定重定向的URL
		c.Request.URL.Path = "/final"
		r.HandleContext(c)
	})
	r.GET("/final", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	r.Run()
}
