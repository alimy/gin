package main

import (
	"net/http"

	"github.com/alimy/gin"
	"github.com/thinkerou/favicon"
)

func main() {
	app := gin.Default()
	_ = favicon.New("./favicon.ico")
	// app.Use(favicon.New("./favicon.ico"))
	app.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello favicon.")
	})
	app.Run(":8080")
}
