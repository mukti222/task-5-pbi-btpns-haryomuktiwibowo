

package controllers

import (
    // "database/sql"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "task-5-pbi-btpns-haryomuktiwibowo/database"
    // "task-5-pbi-btpns-haryomuktiwibowo/models"
    "time"
)

// AddPhoto adalah fungsi untuk menangani endpoint POST /photos
func AddPhoto(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.NotFound(w, r)
        return
    }

    // Parse body form-data
    err := r.ParseMultipartForm(10 << 20) // Limit ukuran file 10 MB
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    // Mendapatkan nilai dari form-data
    title := r.FormValue("title")
    caption := r.FormValue("caption")
    userID := r.FormValue("userID")

    // Mendapatkan file foto dari form-data
    file, handler, err := r.FormFile("photoUrl")
    if err != nil {
        http.Error(w, "Failed to get photo from form data", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Menyimpan file foto ke sistem file server
// Menyimpan file foto ke sistem file server di direktori uploads di dalam direktori proyek
filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
//lokasi file uploads
filepath := filepath.Join("uploads", filename)
f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
if err != nil {
http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
return
}
defer f.Close()

    // Menyalin isi file foto ke sistem file server
    _, err = io.Copy(f, file)
    if err != nil {
        http.Error(w, "Failed to write file to server", http.StatusInternalServerError)
        return
    }

    // Insert foto baru ke database
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO \"Photo\" (Title, Caption, PhotoUrl, UserID) VALUES ($1, $2, $3, $4)", title, caption, filepath, userID)
    if err != nil {
        http.Error(w, "Failed to insert photo into the database", http.StatusInternalServerError)
        return
    }

    // Kirim respons berhasil
    w.WriteHeader(http.StatusCreated)
}
