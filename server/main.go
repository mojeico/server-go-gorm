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
	loggingVideo    repository.VideoLoggingRepository = repository.NewVideoLoggingRepository()
	videoRepository repository.VideoRepository        = repository.NewVideoRepository()
	videoService    service.VideoService              = service.NewVideoService(videoRepository)
	videoController controller.VideoController        = controller.NewVideoController(videoService)
)

func main() {

	server := gin.New()

	server.Use(gin.Recovery(), middlewares.Logger(), middlewares.BasicAuth())

	server.GET("/posts", func(context *gin.Context) {
		loggingVideo.LoggingFindAll("successfully")
		context.JSON(200, videoController.FindAll())
	})

	server.POST("/videos", func(context *gin.Context) {
		if err := videoController.Create(context); err != nil {
			loggingVideo.LoggingCreate("error")
			context.JSON(http.StatusOK, gin.H{"message": "Data is valid"})
		} else {
			loggingVideo.LoggingCreate("successfully")
			context.JSON(http.StatusBadRequest, gin.H{"message": "Data is not valid"})
		}
	})
	server.Run(":8080")

}
