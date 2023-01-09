package main

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/frontend"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func main() {
	gin.SetMode("release")
	config := frontend.LoadConfig()
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Static("/", config.StaticFolder)
	r.NoRoute(func(context *gin.Context) {
		context.File(filepath.Join(config.StaticFolder, "index.html"))
	})
	fmt.Println(config)

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":" + config.InnerPort); err != nil {
		log.Fatal(err)
	}
}
