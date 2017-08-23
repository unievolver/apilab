package main

import (
	"apilab/controllers"
	"apilab/jwtAuth"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {
	defer controllers.DB.Close()
	router := gin.Default()
	router.Use(jwtAuth.JwtCors())
	router.POST("/login", controllers.Login)

	admin := router.Group("/", jwtAuth.JwtMW())
	admin.POST("/permissions", t1)

	//使用域名lab的证书
	router.RunTLS(":443", "./certificate/lab.cert", "./certificate/lab.key")
}

func t1(c *gin.Context) {
	if claims, ok := c.Get("claims"); ok {
		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			c.JSON(http.StatusOK, gin.H{"permissions": mapClaims["permissions"]})
		}
	}
}
