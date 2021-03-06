package controller

import (
	"github.com/AH-dark/random-donate/pkg/encrypt"
	"github.com/AH-dark/random-donate/pkg/response"
	"github.com/AH-dark/random-donate/pkg/utils"
	"github.com/AH-dark/random-donate/service"
	"github.com/gin-gonic/gin"
)

const sessNamespace = "user_info"

// SessionUserHandler 从 Session 获取用户信息
func SessionUserHandler(c *gin.Context) {
	sess := utils.GetSession(c, sessNamespace)
	if sess == nil {
		response.DataHandle(c, nil)
		return
	}

	userId := sess.(uint)
	userInfo, err := service.GetUserById(userId)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	response.DataHandle(c, userInfo)
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	type userLoginData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var loginData userLoginData
	err := c.BindJSON(&loginData)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	// login
	pass := encrypt.Pass(loginData.Password)
	user, err := service.Login(loginData.Login, pass)
	if err != nil {
		response.ServerErrorHandle(c, err)
		return
	}

	// save session
	utils.SetSession(c, map[string]interface{}{
		sessNamespace: user.ID,
	})

	response.DataHandle(c, user)
	return
}
