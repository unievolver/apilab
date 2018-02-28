package controllers

import (
	"apilab/jwtAuth"
	"apilab/models"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//Login 提供给main函数登陆
func Login(c *gin.Context) {
	var user struct {
		UserName string
		Password string
		Checked  bool
	}
	if c.BindJSON(&user) == nil {
		//查找数据库是否存在该用户并比较密码
		var u models.User
		if !db.Where("user_name = ?", user.UserName).First(&u).RecordNotFound() {
			if u.Password == user.Password {
				//找出对应用户的权限
				var permissions []string
				var us []models.User
				us = append(us, u)
				db.Model(&models.Role{}).Related(&us, "Users").Pluck("distinct role_name", &permissions)

				fmt.Println(permissions)
				//设置jwttoken的claims及对应的权限
				claims := jwt.MapClaims{
					"iss":         "evolver",
					"exp":         time.Now().Add(time.Hour * 24 * 7).Unix(),
					"permissions": permissions,
				}
				token, err := jwtAuth.CreateToken(claims)
				if err != nil {
					panic("创建jwttoken失败")
				}
				c.JSON(http.StatusOK, gin.H{
					"code":     200,
					"msg":      "正确",
					"jwttoken": token,
				})
				return
			}
		}
	}
	//否则
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户名或密码错误",
	})
	return
}
