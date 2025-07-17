package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Slug       string         `gorm:"size:200;not null;unique" json:"slug"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Summary    string         `gorm:"size:500" json:"summary"`
	AuthorID   uint           `gorm:"not null" json:"author_id"`
	Author     User           `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID uint           `gorm:"not null" json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Status     string         `gorm:"size:20;not null;default:'draft'" json:"status"` // draft, published
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
