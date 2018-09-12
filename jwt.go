package ginp

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	TokenExpire = 7 * 24 * time.Hour
)

var (
	TokenExpired = errors.New("token expired")
	TokenInvalid = errors.New("invalid token")
)

var (
	emptyTimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
)

type UserInfo struct {
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Avatar string   `json:"avatar"`
	Roles  []string `json:"roles"`
}

type Claims struct {
	*jwt.StandardClaims
	*UserInfo
}

func NewJWT(key string) *JWT {
	return &JWT{
		SignKey: []byte(key),
	}
}

type JWT struct {
	SignKey []byte
}

func (j *JWT) GetKey(token *jwt.Token) (interface{}, error) {
	return j.SignKey, nil
}

func (j *JWT) Create(user *UserInfo) (string, error) {
	claims := &Claims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpire).Unix(),
		},
		UserInfo: user,
	}
	return j.sign(claims)
}

func (j *JWT) Parse(val string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(val, &Claims{}, j.GetKey)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func (j *JWT) Refresh(val string) (string, error) {
	jwt.TimeFunc = emptyTimeFunc
	token, err := jwt.ParseWithClaims(val, &Claims{}, j.GetKey)
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(TokenExpire).Unix()
		return j.sign(claims)
	}
	return "", TokenInvalid
}

func (j *JWT) sign(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

func Auth(obj Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkToken || !obj.TokenRequired(c.Request.URL.Path) {
			return
		}
		token := c.Request.Header.Get("Authorization")
		if parts := strings.Split(token, " "); len(parts) == 2 {
			token = parts[1]
		}
		j := NewJWT(tokenSignKey)
		claims, err := j.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, Response{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}
		c.Set("claims", claims)
	}
}
