// gin框架
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	//用户是查询数据的，用 GET，参数放到查询参数里面，即 ?a=123 这种。
	//用户是提交数据的，用 POST，参数全部放到 Body 里面。
	r.GET("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	//参数路由 http://localhost:8080/user/zzp
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "我是"+name)
	})
	//通配符路由 http://localhost:8080/views/ok.html
	r.GET("/views/*.html", func(c *gin.Context) {
		g := c.Param(".html")
		c.String(http.StatusOK, "页面是%s", g)
	})
	//查询参数 http://localhost:8080/users?id=123
	r.GET("/users", func(c *gin.Context) {
		id := c.Query("id")
		c.String(http.StatusOK, "id的值是%s", id)
	})
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run()
}
