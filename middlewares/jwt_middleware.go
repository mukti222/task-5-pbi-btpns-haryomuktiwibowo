// project/middlewares/jwt_middleware.go

package middlewares

import (
    "net/http"
    "strings"
	"context"
    "task-5-pbi-btpns-haryomuktiwibowo/helpers"
)

// Authenticate adalah middleware untuk memeriksa token JWT pada setiap permintaan
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Mendapatkan token dari header Authorization
        tokenHeader := r.Header.Get("Authorization")
        if tokenHeader == "" {
            http.Error(w, "Missing authorization token", http.StatusUnauthorized)
            return
        }

        // Memisahkan "Bearer" dari token
        parts := strings.Split(tokenHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
            return
        }

        // Verifikasi token JWT
        claims, err := helpers.VerifyToken(parts[1])
        if err != nil {
            http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
            return
        }

        // Menyimpan informasi pengguna dari token ke dalam konteks permintaan
        ctx := context.WithValue(r.Context(), "userID", claims.Email)

        // Memanggil handler berikutnya dengan konteks yang diperbarui
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
