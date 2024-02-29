// project/database/config.go
package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

// ConnectDB berfungsi untuk membuat koneksi ke database PostgreSQL
func ConnectDB() (*sql.DB, error) {
    // Kredensial koneksi ke database
    const (
        host     = "localhost"
        port     = 5432
        user     = "postgres"
        password = "mukti"
        dbname   = "projectbtpns"
    )

    // Membuat string koneksi
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Membuka koneksi ke database
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }

    // Menguji koneksi ke database
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    fmt.Println("Successfully connected to the database!")

    return db, nil
}
