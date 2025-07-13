package models

import (
	"github.com/dgrijalva/jwt-go"
)

// 自定义token结构体
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
