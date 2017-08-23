package jwtAuth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// usage:
// //设置需要的claims
// claims := map[string]interface{}{
// 	"iss": "evolver",
// 	"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
// "permissons": []string{"admin", "backend"},
// }
// jwtClaims := jwtAuth.CovertToJWTClaims(claims)
// fmt.Println(jwtClaims["exp"])
// //获取token
// token, err := jwtAuth.CreateToken(jwtClaims)
// fatal(err)
// fmt.Println(token)
// //获取token里面的claims
// jwtClaims, err = jwtAuth.ParseToken(token)
// fatal(err)
// fmt.Println(jwtClaims["iss"])
// //获取到期剩余时间
// RemainedTime := jwtAuth.GetRemainedTime(jwtClaims)
// fmt.Println(RemainedTime)

var (
	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
)

func init() {
	var err error
	//生成私钥
	privKey, err = rsa.GenerateKey(rand.Reader, 2048)
	fatal(err)
	//获取对应的公钥
	pubKey = &privKey.PublicKey
	//将rsa对象存储在文件里
	// rsaToPem("")
	//将pem文件的字符串转化为rsa对象
	// pemToRsa("")
}

func rsaToPem(filepath string) {
	privKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	filename := path.Join(filepath, "privKey.pem")
	fmt.Println(filename)
	err := ioutil.WriteFile(filename, privKeyBytes, os.ModePerm)
	fatal(err)
	derPkix, err := x509.MarshalPKIXPublicKey(pubKey)
	fatal(err)
	pubKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	})
	filename = path.Join(filepath, "pubKey.pem")
	err = ioutil.WriteFile(filename, pubKeyBytes, os.ModePerm)
	fatal(err)
}

func pemToRsa(filepath string) {
	filename := path.Join(filepath, "privKey.pem")
	privKeyBytes, err := ioutil.ReadFile(filename)
	fatal(err)
	filename = path.Join(filepath, "pubKey.pem")
	pubKeyBytes, err := ioutil.ReadFile(filename)
	fatal(err)
	privKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	fatal(err)
	pubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	fatal(err)
}

func CreateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := token.SignedString(privKey)
	return tokenString, err
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	//token解析失败
	if token == nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func GetRemainedTime(claims jwt.MapClaims) time.Duration {
	//直接jwt.MapClaims定义的为float64,通过CovertToJWTClaims过来的exp是int64类型
	if expFloat, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(expFloat), 0)
		remainedTime := expTime.Sub(time.Now())
		return remainedTime
	}
	if expInt64, ok := claims["exp"].(int64); ok {
		expTime := time.Unix(expInt64, 0)
		remainedTime := expTime.Sub(time.Now())
		return remainedTime
	}
	return 0
}

//从客户端获取jwt token
func getJwt(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("auth header empty")
	}
	authHeaderArray := strings.SplitN(authHeader, " ", 2)
	if len(authHeaderArray) == 2 && authHeaderArray[0] == "Bearer" {
		return strings.TrimSpace(authHeaderArray[1]), nil
	}
	return "", errors.New("invalid authorization header")
}

func JwtMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从客户端获取jwt token
		tokenstring, err := getJwt(c)
		if err != nil || tokenstring == "" {
			//token为空
			c.JSON(http.StatusUnauthorized, gin.H{"status": "-1", "msg": err.Error()})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenstring)
		if err != nil {
			//token无效
			c.JSON(http.StatusUnauthorized, gin.H{"status": "-2", "msg": err.Error()})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}

}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JwtCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	// 允许头部文件使用字段Authorization
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	//允许跨域访问
	//生产服务器请注释或使用config.AllowOrigins
	// config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowAllOrigins = true
	return cors.New(config)
}
