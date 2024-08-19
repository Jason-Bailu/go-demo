# Go web Gin

- Gin安装：**go get -u github.com/gin-gonic/gin**

- Gin文档：[Gin DOC](https://gin-gonic.com/zh-cn/)
  - Gin原理：Gin框架中的路由使用的是[httprouter](https://github.com/julienschmidt/httprouter)这个库。

## 快速入门

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
```

## RESTfulAPI

```go
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
```

## 参数读取和返回

```go
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
```

## 文件上传

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	file_upload_path := "./files"

	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	// 单个文件
	r.POST("/file/upload", func(c *gin.Context) {
		// 读取文件form
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Print(file.Filename)
		// 保存文件
		dst := fmt.Sprintf(file_upload_path+"/%s", file.Filename)
		c.SaveUploadedFile(file, dst)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})

	// 多个文件
	r.POST("/files/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file"]
		for _, file := range files {
			log.Println(file.Filename)
			dst := fmt.Sprintf(file_upload_path+"/%s", file.Filename)
			c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})

	r.Run()
}
```

## 重定向

```go
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
```

## 路由组

```go
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
```

## 中间件扩展



## 多个项目同时启动

- go get -u golang.org/x/sync/errgroup

```go
package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

var (
	g errgroup.Group
)

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})
	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 02",
			},
		)
	})
	return e
}

func main() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// 借助errgroup.Group或者自行开启两个goroutine分别启动两个服务
	g.Go(func() error {
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
```

