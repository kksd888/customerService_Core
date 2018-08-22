package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{"code": "ok"})
	})

	r.Run(":5000")
}
