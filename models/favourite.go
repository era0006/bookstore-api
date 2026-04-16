package models

import "time"

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	BookID    uint      `gorm:"not null" json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}
