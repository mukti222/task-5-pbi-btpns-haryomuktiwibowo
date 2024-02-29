package main

import (
    "log"
    "net/http"
    "task-5-pbi-btpns-haryomuktiwibowo/database"
    "task-5-pbi-btpns-haryomuktiwibowo/router"
)

func main() {
    // Menghubungkan ke database
    db, err := database.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }
    defer db.Close()

    // Menginisialisasi router dan menetapkan endpoint ke controller yang sesuai
    router.InitRouter()

    // Menjalankan server HTTP di localhost:8080
    log.Println("Server started on localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
