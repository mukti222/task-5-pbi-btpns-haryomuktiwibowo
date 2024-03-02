// project/models/photo_model.go

package models

import "mime/multipart"

// Photo adalah struktur untuk merepresentasikan data foto
type Photo struct {
    ID        int            `json:"id"`
    Title     string         `json:"title"`
    Caption   string         `json:"caption"`
    PhotoUrl *multipart.File `json:"-"`
    UserID    int            `json:"userID"`
}
