package model

import "time"

// News struct is entity for news table
type News struct {
	ID      uint      `gorm:"primary_key" json:"id"`
	Author  string    `gorm:"type:text" json:"author"`
	Body    string    `gorm:"type:text" json:"body"`
	Created time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created"`
}
