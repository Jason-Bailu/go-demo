package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// get
	r.GET("/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GET",
		})
	})

	// post
	r.POST("/post", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POST",
		})
	})

	// put
	r.PUT("/put", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PUT",
		})
	})

	// DELETE
	r.DELETE("/delete", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DELETE",
		})
	})

	r.Run()
}
