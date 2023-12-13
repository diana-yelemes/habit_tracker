// models/models.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model
	ID                  uint   `json:"id"`
	Habit_Name          string `json:"habit_name"`
	CalendarDays        []CalendarDay
	Target_Repeat_Count int        `json:"target_repeat_count"`
	Repeat_Count        int        `json:"repeat_count"`
	Notes               string     `json:"notes"`
	Completed           bool       `json:"completed"`
	CompletionDate      *time.Time `json:"completion_date"`
}

type CalendarDay struct {
	gorm.Model
	DayOfWeek string `gorm:"column:day_of_week"`
	Completed bool
	HabitID   uint
}
