package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	group1 := r.Group("/group1")
	{
		group1.GET("/index1", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"messgae": "/group1/index1",
			})
		})
		group1.GET("/index2", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"messgae": "/group1/index2",
			})
		})
	}

	group2 := r.Group("/group2")
	{
		group2.GET("/index1", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"messgae": "/group2/index1",
			})
		})
		group2.GET("/index2", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"messgae": "/group2/index2",
			})
		})
	}

	r.Run()
}
