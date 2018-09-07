package ginp

import "time"

type Model struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"create_at" gorm:"index;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"update_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
