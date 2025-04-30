package jwt_pkg

import "github.com/golang-jwt/jwt/v5"

type JwtProducer interface {
	Produce(key, userId string, startTime string, deviceHash string) string
}

var _JwtProducer = (*JwtProducerImpl)(nil)

type JwtProducerImpl struct {
}

func (j *JwtProducerImpl) Produce(key, userId string, startTime string, deviceHash string) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"iat":     startTime,
		"device":  deviceHash,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	return tokenString
}
