// project/app/user.go
package app

// User adalah struktur untuk merepresentasikan data pengguna
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
