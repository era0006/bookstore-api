package models

type Book struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Title      string  `gorm:"not null" json:"title"`
	AuthorID   uint    `json:"author_id"`
	CategoryID uint    `json:"category_id"`
	Price      float64 `json:"price"`
}
