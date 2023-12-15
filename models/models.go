// models/models.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	RepeatCount       int        `json:"repeat_count"`
	TargetRepeatCount int        `json:"target_repeat count"`
	Monday            bool       `json:"monday"`
	Tuesday           bool       `json:"tuesday"`
	Wednesday         bool       `json:"wednesday"`
	Thursday          bool       `json:"thursday"`
	Friday            bool       `json:"friday"`
	Saturday          bool       `json:"saturday"`
	Sunday            bool       `json:"sunday"`
	Completed         bool       `json:"completed"`
	CompletionDate    *time.Time `json:"completion_date"`
}
