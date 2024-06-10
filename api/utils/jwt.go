package utils

import (
	"api/config"
	"api/global"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JWT struct {
	JwtKey []byte
}

func NewJWT(key string) *JWT {
	return &JWT{
		[]byte(key), // 密钥加密
	}
}

type MyClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

// 定义错误
var (
	TokenExpired     = errors.New("过期token")
	TokenNotValidYet = errors.New("无效token")
	TokenMalformed   = errors.New("错误token")
	TokenInvalid     = errors.New("无效token")
)

// CreateToken 生成token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

// ParserToken 解析token
func (j *JWT) ParserToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}

	return nil, TokenInvalid
}

// JwtToken jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    global.ErrCodeJWTMissingToken,
				"message": "Please login in first",
			})
			c.Abort()
			return
		}
		j := NewJWT(config.Conf.JwtSecret)
		// 解析token
		claims, err := j.ParserToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.ErrCodeJWTInvalidToken,
				"message": "Please login in first",
			})
			c.Abort()
			return
		}
		c.Set("use_id", claims.UserId)
		c.Next()
	}
}
