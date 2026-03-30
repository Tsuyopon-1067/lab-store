package model

import "time"

type Backup struct {
	ID        int       `json:"id"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
}
