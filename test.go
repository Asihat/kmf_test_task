package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		c.JSON(200, "hey its post request")
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hey its get request")
	})
	r.Run(":8081")
}
