package request

import (
	"github.com/dgrijalva/jwt-go"
)

// Custom claims structure
type CustomClaims struct {
	ID          uint
	Phone    string
	RealName    string
	RoleId    uint
	BufferTime  int64
	jwt.StandardClaims
}
