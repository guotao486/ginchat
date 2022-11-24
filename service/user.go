/*
 * @Author: GG
 * @Date: 2022-11-21 16:46:00
 * @LastEditTime: 2022-11-23 17:34:40
 * @LastEditors: GG
 * @Description:
 * @FilePath: \ginchat\service\user.go
 *
 */
package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Login
// @Summary login
// @Schemes user login
// @Decription user login
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @param name formData string true "用户名"
// @param password formData string true "新密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/login [post]
func Login(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	currentUser := models.FindUserByName(name)
	if currentUser.Name == "" {
		c.JSON(-1, gin.H{
			"message": "不存在该用户",
		})
		return
	}

	flag := utils.ValidPassword(password, currentUser.Salt, currentUser.PassWord)
	if !flag {
		c.JSON(-1, gin.H{
			"message": "密码不正确",
		})
		return
	}

	timeStr := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(timeStr)
	utils.DB.Model(&currentUser).Where("id = ?", currentUser.ID).Update("identity", temp)
	c.JSON(200, gin.H{
		"data": currentUser,
	})
}

// GetUserList
// @Summary 获取用户列表
// @Schemes
// @Decription get user list
// @Tags user
// @Accept json
// @Produce json
// @Success 200 UserBasic json{"code":200, "data":[]models.user_basic}
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	list := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// CreateUser
// @Summary 创建用户
// @Schemes create user
// @Decription createUser
// @Tags user
// @Accept json
// @Produce json
// @param name query string false "用户名"
// @param password query string false "用户密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/create [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")

	currentUser := models.FindUserByName(user.Name)

	if currentUser.Name != "" {
		c.JSON(-1, gin.H{
			"message": "用户名已注册",
		})
		return
	}

	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致！",
		})
		return
	}
	user.Salt = fmt.Sprintf("%06d", rand.Int31())
	user.PassWord = utils.MakePassword(password, user.Salt)

	models.Create(&user)
	c.JSON(200, gin.H{
		"message": "Create User Success",
	})
}

// DeleteUser
// @Summary 删除用户
// @Schemes delete user
// @Decription deleteUser
// @Tags user
// @Accept json
// @Produce json
// @param id query int true "用户id"
// @Success 200 {string} json{"code","message"}
// @Router /user/delete [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.Delete(&user)
	c.JSON(200, gin.H{
		"message": "delete user success",
	})
}

// UpdateUser
// @Summary 修改用户
// @Schemes update user
// @Decription updateUser
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @param id formData int true "用户id"
// @param name formData string true "用户名"
// @param password formData string true "新密码"
// @param repassword formData string true "旧密码"
// @param phone formData string true "手机号"
// @param email formData string true "邮箱"
// @Success 200 {string} json{"code","message"}
// @Router /user/update [patch]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err,
		})
		return
	}

	if password == repassword {
		c.JSON(-1, gin.H{
			"message": "password and oldpassword 不能一样",
		})
		return
	}

	user.PassWord = password
	result := models.Update(&user)
	if result.Error != nil {
		c.JSON(-1, gin.H{
			"message": result.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "update user success",
	})
}

//防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)

	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}

}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
