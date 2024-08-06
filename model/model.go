package model

import (
    "time"
)
type Material struct {
    UUID         string    `json:"uuid"`
    MaterialType string    `json:"material_type"`
    Status       string    `json:"status"`
    Title        string    `json:"title"`
    Content      string    `json:"content"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
