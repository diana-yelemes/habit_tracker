package models

import (
	"time"

	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model
	ID                  uint       `json:"id"`
	Habit_Name          string     `json:"habit_name"`
	Target_Repeat_Count uint       `json:"target_repeat_count"`
	Repeat_Count        uint       `json:"repeat_count"`
	Notes               string     `json:"notes"`
	Completed           bool       `json:"completed"`
	CompletionDate      *time.Time `json:"completion_date"`
}
