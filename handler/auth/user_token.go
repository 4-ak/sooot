package authtoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	TokenFormatNotMatchError = "UserToken : couldn't found jwtMap"
	NotExistParameter        = "UserToken : %v not exists in jwt"
)

type UserToken struct {
	ID         string
	ClientID   string
	Permission string
	Expires    string
}

func NewUserToken(jwtCookie string) (UserToken, error) {
	token, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return UserToken{}, err
	}

	jwtMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserToken{}, fmt.Errorf(TokenFormatNotMatchError)
	}

	var user UserToken

	//TODO : think elegant way convert map to struct
	user.ID, ok = jwtMap["id"].(string)
	if !ok {
		return UserToken{}, fmt.Errorf(NotExistParameter, "ID")
	}

	user.ClientID, ok = jwtMap["cid"].(string)
	if !ok {
		return UserToken{}, fmt.Errorf(NotExistParameter, "ClientID")
	}

	user.Permission, ok = jwtMap["per"].(string)
	if !ok {
		return UserToken{}, fmt.Errorf(NotExistParameter, "Permission")
	}

	user.Expires, ok = jwtMap["expires"].(string)
	if !ok {
		return UserToken{}, fmt.Errorf(NotExistParameter, "Expires")
	}
	return user, nil
}

func (ut *UserToken) JWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      ut.ID,
		"cid":     ut.ClientID,
		"per":     ut.Permission,
		"expires": ut.Expires,
	})
	return token.SignedString([]byte("secret"))
}

func (ut *UserToken) isExpiration() bool {
	expires, err := time.Parse(time.RFC3339, ut.Expires)
	if err != nil {
		return false
	}
	if expires.Before(time.Now()) {
		return false
	}
	return true
}
