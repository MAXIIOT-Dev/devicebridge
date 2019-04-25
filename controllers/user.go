package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxiiot/devicebridge/storage"
	"github.com/maxiiot/devicebridge/utils"
)

// User login info
type User struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// UserPassword user change password info
type UserPassword struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}

// Login 登陆
// @summary 登陆
// @description 用户登陆
// @tags user
// @accept json
// @produce json
// @param user body controllers.User true "user login info"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @router /user/login [post]
func Login(c *gin.Context) {
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "请输入用户名和密码", nil)
		return
	}

	usr, err := storage.LoginUser(user.UserName)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "用户名不存在", nil)
		return
	}

	err = utils.Compare(user.Password, usr.PasswordHash)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "密码错误", nil)
		return
	}

	token, err := CreateToken(usr)
	if err != nil {
		msg := fmt.Sprintf("生成JWT错误: %s", err)
		Response(c, http.StatusInternalServerError, 1, msg, nil)
		return
	}

	Response(c, http.StatusOK, 0, "success", gin.H{
		"jwt": token,
	})
}

// @summary 更改密码
// @description 更改密码
// @tags user
// @accept json
// @produce json
// @param user body controllers.UserPassword true "user password info"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @router /user/changepwd [put]
// @security ApiKeyAuth
func ChangeUserPassword(c *gin.Context) {
	var user_pwd UserPassword
	err := c.ShouldBind(&user_pwd)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	username := c.GetString("username")

	olduser, err := storage.LoginUser(username)
	if err != nil {
		Response(c, http.StatusUnauthorized, 1, "请先登陆", nil)
		return
	}

	err = utils.Compare(user_pwd.OldPassword, olduser.PasswordHash)
	if err != nil {
		Response(c, http.StatusUnauthorized, 1, "旧密码错误", nil)
		return
	}

	olduser.PasswordHash, err = utils.Hash(user_pwd.NewPassword)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	err = storage.UpdateUserPassword(olduser)
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusOK, 0, "success", nil)
}

// CreateUser create user
// @summary 新增用户
// @description 新增用户
// @tags user
// @accept json
// @produce json
// @param user body controllers.User true "create user info"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /user/add [post]
func CreateUser(c *gin.Context) {
	var usr User
	err := c.ShouldBind(&usr)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "用户名或密码不能为空", nil)
		return
	}

	pwdHash, err := utils.Hash(usr.Password)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "用户名或密码不能为空", nil)
		return
	}

	err = storage.CreateUser(storage.User{
		UserName:     usr.UserName,
		PasswordHash: pwdHash,
	})
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusBadRequest, 0, "success", nil)
}
