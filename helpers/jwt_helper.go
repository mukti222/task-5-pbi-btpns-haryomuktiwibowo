// project/helpers/jwt_helper.go
package helpers

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

// Claims adalah struktur yang digunakan untuk menetapkan klaim JWT
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// GenerateToken menghasilkan token JWT dengan email pengguna
func GenerateToken(email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Token akan berlaku selama 24 jam

    claims := &Claims{
        Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// VerifyToken memverifikasi token JWT
func VerifyToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, err
    }

    return claims, nil
}
