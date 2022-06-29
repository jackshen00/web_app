package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

//const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("asasasasa")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userid int64, username string) (string, error) {
	c := MyClaims{
		userid,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),
			Issuer:    "my-project",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	//claims, ok := token.Claims.(*MyClaims)
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
