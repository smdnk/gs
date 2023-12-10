package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type T struct {
	PkgId string `form:"pkgId" json:"pkgId" xml:"pkgId" binging:"required"`
}

func (t T) toString() {
	log.Println("-----------------------")
}

func setStatic(r *gin.Engine) {

	r.Static("/assets", "web/assets")
	r.LoadHTMLGlob("web/page/*")
}

func MyLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 这里也是我们自己的组件
		fmt.Println("自己的组件")
	}
}

// Default 这里就是gin框架的实现
func Default() *gin.Engine {
	engine := gin.New()
	// gin默认挂载了logger和recovery组件
	engine.Use(gin.Logger(), gin.Recovery())
	// 我们也可以挂载自己的中间件
	engine.Use(func(c *gin.Context) {
		// 这里就是我们自己的组件
		fmt.Println("自己的组件")
		//	 这里面可以做一些拦截，处理一些逻辑
	})
	// 我们也可以仿官方写
	engine.Use(MyLogger())

	return engine
}

func main() {

	//r := gin.Default()
	r := Default()
	setStatic(r)

	r.GET("index", func(context *gin.Context) {
		context.HTML(200, "index.html", nil)
	})

	// 127.0.0.1/index/999?www=444
	r.POST("/index/:ppp", func(c *gin.Context) {

		param := c.Param("ppp")   // 999 // 获取路由占位符号对应的数据
		value := c.Query("www")   // 444 获取请求路径拼接的数据
		form := c.PostForm("pkg") // 获取form表单的请求数据
		params := c.Params

		var t T // c.ShouldBindJSON(&t) 解析 JSON 数据
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Param()%v,query() %v;postForm()%v;params%v;jsonPost:%v", param, value, form, params, t)

		c.JSONP(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/someJson", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}
		c.AsciiJSON(http.StatusOK, data)
	})

	err := r.Run(":9090")
	if err != nil {
		log.Fatalln("服务器启动失败")
		return
	}

}
