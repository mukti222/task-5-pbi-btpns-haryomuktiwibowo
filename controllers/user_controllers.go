// project/controllers/user_controllers.go
package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "strconv"
    "task-5-pbi-btpns-haryomuktiwibowo/database"
    "task-5-pbi-btpns-haryomuktiwibowo/helpers"
    "golang.org/x/crypto/bcrypt"
)

// LoginPayload adalah struktur untuk merepresentasikan payload saat login pengguna
type LoginPayload struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
// LoginUser adalah fungsi untuk menangani endpoint POST /users/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.NotFound(w, r)
        return
    }

    // Mendekode payload JSON
    var payload LoginPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Mendapatkan pengguna berdasarkan email dari database
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var storedPassword string
    err = db.QueryRow("SELECT Password FROM \"User\" WHERE Email = $1", payload.Email).Scan(&storedPassword)
    if err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Memverifikasi kata sandi
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(payload.Password)); err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Kata sandi benar, generate token JWT
    token, err := helpers.GenerateToken(payload.Email)
    if err != nil {
        http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
        return
    }

    // Kirim token JWT sebagai respons
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RegisterPayload adalah struktur untuk merepresentasikan payload saat registrasi pengguna
type RegisterPayload struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// RegisterUser adalah fungsi untuk menangani endpoint POST /users/register
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.NotFound(w, r)
        return
    }

    // Mendekode payload JSON
    var payload RegisterPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Validasi payload
    if payload.Username == "" || payload.Email == "" || payload.Password == "" {
        http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
        return
    }

    // Hash password menggunakan bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    // Mendapatkan waktu saat ini
    currentTime := time.Now()

    // Insert pengguna baru ke database
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO \"User\" (Username, Email, Password, Created_At, Updated_At) VALUES ($1, $2, $3, $4, $5)", payload.Username, payload.Email, hashedPassword, currentTime, currentTime)
    if err != nil {
        http.Error(w, "Failed to insert user", http.StatusInternalServerError)
        return
    }

    // Kirim respons berhasil
    w.WriteHeader(http.StatusCreated)
}


// UpdatePayload adalah struktur untuk merepresentasikan payload saat update pengguna
type UpdatePayload struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// DeleteUser adalah fungsi untuk menangani endpoint DELETE /users/:userId
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.NotFound(w, r)
        return
    }

    // Mendapatkan userID dari URL
    userID, err := strconv.Atoi(r.URL.Path[len("/users/"):])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Mendapatkan pengguna berdasarkan userID dari database
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM \"User\" WHERE ID = $1", userID)
    if err != nil {
        http.Error(w, "Failed to delete user", http.StatusInternalServerError)
        return
    }

    // Kirim respons berhasil
    w.WriteHeader(http.StatusOK)
}

// UpdateUser adalah fungsi untuk menangani endpoint PUT /users/:userId
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.NotFound(w, r)
        return
    }

    // Mendapatkan userID dari URL
    userID, err := strconv.Atoi(r.URL.Path[len("/users/"):])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Mendekode payload JSON
    var payload UpdatePayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Validasi payload
    if payload.Username == "" || payload.Email == "" || payload.Password == "" {
        http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
        return
    }

    // Hash password menggunakan bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    // Mendapatkan waktu saat ini
    currentTime := time.Now()

    // Update pengguna ke database
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("UPDATE \"User\" SET Username = $1, Email = $2, Password = $3, Updated_At = $4 WHERE ID = $5", payload.Username, payload.Email, hashedPassword, currentTime, userID)
    if err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    // Kirim respons berhasil
    w.WriteHeader(http.StatusOK)
}