package main

import (
	"github.com/gin-gonic/gin"
	"myMod/controller"
	"myMod/middlewares"
	"myMod/repository"
	"myMod/service"
	"net/http"
)

var (
	videoRepository repository.VideoRepository = repository.NewVideoRepository()
	videoService    service.VideoService       = service.New(videoRepository)
	videoController controller.VideoController = controller.New(videoService)
)

func main() {

	server := gin.New()

	server.Use(gin.Recovery(), middlewares.Logger(), middlewares.BasicAuth())

	server.GET("/posts", func(context *gin.Context) {
		context.JSON(200, videoController.FindAll())
	})

	server.POST("/videos", func(context *gin.Context) {
		if err := videoController.Create(context); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Data is not valid"})
		} else {
			context.JSON(http.StatusOK, gin.H{"message": "Data is valid"})
		}
	})
	server.Run(":8080")

}
