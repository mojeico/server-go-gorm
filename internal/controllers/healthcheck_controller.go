package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/pkg/logger"
)

func HealthCheck(ctx *gin.Context) {

	res := response.BuildResponse("OK", "service is working")

	ctx.JSON(http.StatusOK, res)

	logger.SystemLoggerInfo("service is working").Info("OK")

}
