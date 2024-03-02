

package controllers

import (
    // "database/sql"
    // "encoding/json"
    "fmt"
    "io"
    "net/http"
    "io/ioutil"
    "os"
    "strconv"
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


func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
    // Mendapatkan ID foto dari URL
    photoID := r.URL.Query().Get("photoId")

    // Mendapatkan form-data yang dikirim
    err := r.ParseMultipartForm(10 << 20) // Menggunakan 10MB sebagai batas maksimum ukuran file
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    // Mengambil file foto dari form-data
    file, handler, err := r.FormFile("photoUrl")
    if err != nil {
        http.Error(w, "Failed to get photo from form data", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Membaca data file
    data, err := ioutil.ReadAll(file)
    if err != nil {
        http.Error(w, "Failed to read file data", http.StatusInternalServerError)
        return
    }

    // Mendapatkan ekstensi file
    extension := filepath.Ext(handler.Filename)

    // Menyimpan file di sistem
    // Menghasilkan nama unik untuk file
    fileName := fmt.Sprintf("%s%s", photoID, extension)
    // Menyimpan file di folder "uploads"
    fileLocation := filepath.Join("uploads", fileName)
    err = ioutil.WriteFile(fileLocation, data, 0666)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    // Mendapatkan URL file yang disimpan
    // Dalam contoh ini, kita menggunakan URL relatif untuk file yang disimpan
    // Anda dapat menyesuaikan dengan kebutuhan Anda, misalnya menggunakan URL absolut.
    fileURL := "/" + fileLocation

    // Mengupdate URL foto di database
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Jalankan kueri update di database
    _, err = db.Exec("UPDATE \"Photo\" SET PhotoUrl = $1 WHERE ID = $2", fileURL, photoID)
    if err != nil {
        http.Error(w, "Failed to update photo", http.StatusInternalServerError)
        return
    }

    // Kirim respons berhasil
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Photo updated successfully")
}

// DeletePhoto adalah fungsi untuk menangani endpoint DELETE /photos/:photoId
func DeletePhoto(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.NotFound(w, r)
        return
    }

    // Mendapatkan photoId dari URL
    // pakai PARAMS photoId
    photoID := r.URL.Query().Get("photoId")

    fmt.Println("isi:",photoID)

    photoIDInt, err := strconv.Atoi(photoID)
if err != nil {
    // Handle kesalahan konversi
    http.Error(w, "Invalid photo ID", http.StatusBadRequest)
    return
}

    // Tambahkan logika untuk menghapus foto berdasarkan photoId di sini
    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

// Misalnya, Anda dapat menggunakan exec untuk menghapus foto
result, err := db.Exec("DELETE FROM \"Photo\" WHERE ID = $1", photoIDInt)
if err != nil {
    // Tambahkan logging untuk mengetahui lebih detail tentang kesalahan yang terjadi
    fmt.Println("Failed to execute delete query:", err)
    http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
    return
}

// Periksa apakah query berhasil dijalankan
rowsAffected, err := result.RowsAffected()
if err != nil {
    fmt.Println("Failed to get rows affected:", err)
    http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
    return
}

if rowsAffected == 0 {
    // Jika tidak ada baris yang terpengaruh, foto dengan ID yang diberikan mungkin tidak ditemukan
    http.Error(w, "Photo not found", http.StatusNotFound)
    return
}

// Kirim respons berhasil
w.WriteHeader(http.StatusOK)
}