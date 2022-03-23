package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"uranus/common/ctxdata"
)

func GetJwtToken(secretKey string, iat, seconds, id int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[ctxdata.CtxKeyJwtUserId] = id
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
