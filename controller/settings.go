package controller

import (
	"github.com/AH-dark/random-donate/dataType"
	"github.com/AH-dark/random-donate/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetBasicSettings 获取基本站点信息
func GetBasicSettings(c *gin.Context) {
	settings := model.GetSettingByType([]string{"basic"})

	c.JSON(http.StatusOK, &dataType.ApiResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    settings,
	})
}
