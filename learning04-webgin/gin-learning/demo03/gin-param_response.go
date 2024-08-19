package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 参数绑定
type Login struct {
	Account  string `form:"account" json:"account" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	r := gin.Default()

	// params
	r.GET("/get/params", func(c *gin.Context) {
		// 含有默认值
		p1 := c.DefaultQuery("p1", "param1")
		// 不含有默认值
		p2 := c.Query("p2")
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"p1":      p1,
			"p2":      p2,
		})
	})

	// form
	r.POST("/post/form", func(c *gin.Context) {
		// form表单
		item1 := c.PostForm("item1")
		item2 := c.PostForm("item2")
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"item1":   item1,
			"item2":   item2,
		})
	})

	// body
	r.POST("/post/body", func(c *gin.Context) {
		// jsonbody
		body, _ := c.GetRawData()
		var m map[string]interface{}
		// 序列化
		_ = json.Unmarshal(body, &m)
		c.JSON(http.StatusOK, m)
	})

	// path
	r.GET("/get/path/:p1/:p2", func(c *gin.Context) {
		p1 := c.Param("p1")
		p2 := c.Param("p2")
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"p1":      p1,
			"p2":      p2,
		})
	})

	// 参数绑定binding
	// 顺序解析
	// 如果是 GET 请求，只使用 Form 绑定引擎（query）
	// 如果是 POST 请求，首先检查 content-type 是否为 JSON 或 XML，然后再使用 Form（form-data）
	// binding json
	r.POST("/bind/loginJSON", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err == nil {
			fmt.Printf("login info:%#v\n", login)
			c.JSON(http.StatusOK, gin.H{
				"account":  login.Account,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// binding form
	r.POST("/bind/loginForm", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"account":  login.Account,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// binding params
	r.GET("/bind/loginParam", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"account":  login.Account,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.Run()
}
