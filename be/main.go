package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.Static("/uploads", "./uploads")
}
