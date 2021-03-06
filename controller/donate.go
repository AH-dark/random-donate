package controller

import (
	"github.com/AH-dark/random-donate/dataType"
	"github.com/AH-dark/random-donate/dataType/payment"
	"github.com/AH-dark/random-donate/model"
	"github.com/AH-dark/random-donate/pkg/hash"
	"github.com/AH-dark/random-donate/pkg/response"
	"github.com/AH-dark/random-donate/pkg/utils"
	"github.com/AH-dark/random-donate/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"strings"
)

// DonatePostHandler 新增信息
func DonatePostHandler(c *gin.Context) {
	type donateInfoReq struct {
		Name    string                `form:"name" json:"name" validate:"required"`
		Comment string                `form:"comment" json:"comment" validate:"required"`
		QRCode  *multipart.FileHeader `form:"qrcode" json:"-" validate:"required"`
	}

	var data donateInfoReq
	err := c.ShouldBind(&data)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	url, err := utils.ParseQRCode(data.QRCode)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	// generate data
	user, _ := service.GetUserBySession(c)
	dbData := model.DonateInfo{
		Name:    data.Name,
		Comment: data.Comment,
		Payment: "",
		Url:     url,
		Author:  0,
	}
	if user != nil {
		dbData.Author = user.ID
	}

	// check payment and url
	switch true {
	case strings.HasPrefix(url, "https://qr.alipay.com/"):
		dbData.Payment = payment.Alipay
	case strings.HasPrefix(url, "wxp://"):
		dbData.Payment = payment.Wechat
	default:
		c.JSON(http.StatusBadRequest, &dataType.ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "payment is not legal",
		})
		return
	}

	// find if exist
	isExist, err := service.DonateInfoIsExist(&model.DonateInfo{Url: dbData.Url})
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	if isExist {
		c.JSON(http.StatusBadRequest, &dataType.ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "data is exist",
		})
		return
	}

	// insert to database
	err = service.DonateInfoSave(&dbData)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	// get full data
	dbData, err = service.DonateInfoFind(&dbData)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	response.DataHandle(c, hash.Id(dbData.ID, hash.DonateId))
}

// DonateRandomGetHandler 随机获取一条信息
func DonateRandomGetHandler(c *gin.Context) {
	var prevId uint = 0

	// 取值处理
	sessPrevId := utils.GetSession(c, "random_donate_prev_id")
	if sessPrevId != nil {
		prevId = sessPrevId.(uint)
	}

	// 获取数据
	data, err := service.DonateInfoRandomGet(prevId)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	utils.SetSession(c, map[string]interface{}{
		"random_donate_prev_id": data.ID,
	})

	utils.Log().Debug("random donate info: prev: %d, new: %d", sessPrevId, data.ID)

	c.JSON(http.StatusOK, &dataType.ApiResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    hash.Id(data.ID, hash.DonateId),
	})
}

func DonateHashGetHandler(c *gin.Context) {
	hashCode := c.Param("hash")
	id, err := hash.DecodeID(hashCode, hash.DonateId)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	donate, err := service.DonateInfoFind(&model.DonateInfo{
		Model: gorm.Model{
			ID: id,
		},
	})
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	response.DataHandle(c, donate)
}
