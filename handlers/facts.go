package handlers

import (
	"fmt"
	"time"

	"github.com/diana-yelemes/habit_tracker/database"
	"github.com/diana-yelemes/habit_tracker/models"
	"github.com/gofiber/fiber/v2"
)

// Home handles the home route
func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the Habit Tracker!")
}

// GetAllUserHabits retrieves all user habits
func GetAllUserHabits(c *fiber.Ctx) error {
	habits := []models.Habit{}
	database.DB.Db.Find(&habits)

	return c.Status(200).JSON(habits)
}

// CreateHabit adds a new habit
func CreateHabit(c *fiber.Ctx) error {
	habit := new(models.Habit)
	if err := c.BodyParser(habit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	habit.Completed = false
	habit.CompletionDate = nil

	database.DB.Db.Create(&habit)

	return c.Status(200).JSON(habit)
}

// UpdateHabit updates an existing habit
func UpdateHabit(c *fiber.Ctx) error {
	habitID := c.Params("id")
	var existingHabit models.Habit

	if err := database.DB.Db.First(&existingHabit, habitID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Habit not found",
		})
	}

	updatedHabit := new(models.Habit)
	if err := c.BodyParser(updatedHabit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Update the existing habit with the new information
	existingHabit.Habit_Name = updatedHabit.Habit_Name
	existingHabit.Target_Repeat_Count = updatedHabit.Target_Repeat_Count
	existingHabit.Repeat_Count = updatedHabit.Repeat_Count
	existingHabit.Notes = updatedHabit.Notes

	// Save the updated habit to the database
	if err := database.DB.Db.Save(&existingHabit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Habit updated successfully",
		"habit":   existingHabit,
	})
}

// DeleteHabit deletes an existing habit
func DeleteHabit(c *fiber.Ctx) error {
	habitID := c.Params("id")
	var habit models.Habit

	if err := database.DB.Db.First(&habit, habitID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Habit not found",
		})
	}

	database.DB.Db.Delete(&habit)

	return c.Status(200).JSON(fiber.Map{
		"message": "Habit deleted successfully",
	})
}

// GetHabitByID retrieves a single habit by ID
func GetHabitByID(c *fiber.Ctx) error {
	habitID := c.Params("id")
	var habit models.Habit

	if err := database.DB.Db.First(&habit, habitID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Habit not found",
		})
	}

	return c.Status(200).JSON(habit)
}

// FilterHabits filters habits based on query parameters
func FilterHabits(c *fiber.Ctx) error {
	// Extract query parameters from the request
	name := c.Query("name")

	// Build the query string
	var sql string
	sql = "SELECT * FROM habits WHERE deleted_at IS NULL"

	// Append filters based on parameters
	if name != "" {
		sql += fmt.Sprintf(" AND habit_name ILIKE '%%%s%%'", name)
	}

	sql += " ORDER BY id"

	// Execute the dynamic query
	var habits []models.Habit
	if err := database.DB.Db.Raw(sql).Scan(&habits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(habits)
}

// CompleteHabit marks a habit as completed
func CompleteHabit(c *fiber.Ctx) error {
	habitID := c.Params("id")
	var habit models.Habit

	if err := database.DB.Db.First(&habit, habitID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Habit not found",
		})
	}

	// Perform the completion action (update the database, e.g., set completion flag, update completion date)
	habit.Completed = true
	habit.CompletionDate = &[]time.Time{time.Now()}[0]

	if err := database.DB.Db.Save(&habit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Habit marked as completed",
	})
}

// UndoCompleteHabit undoes the completion status of a habit
func UndoCompleteHabit(c *fiber.Ctx) error {
	habitID := c.Params("id")
	var habit models.Habit

	if err := database.DB.Db.First(&habit, habitID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Habit not found",
		})
	}

	// Perform the undo completion action (update the database, e.g., reset completion flag)
	habit.Completed = false
	habit.CompletionDate = nil

	if err := database.DB.Db.Save(&habit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Habit completion status undone",
	})
}

// GetCompletedHabits retrieves all completed habits
func GetCompletedHabits(c *fiber.Ctx) error {
	// Query the database to retrieve completed habits
	var completedHabits []models.Habit
	if err := database.DB.Db.Where("completed = ?", true).Find(&completedHabits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(completedHabits)
}

// GetIncompleteHabits retrieves all incomplete habits
func GetIncompleteHabits(c *fiber.Ctx) error {
	// Query the database to retrieve incomplete habits
	var incompleteHabits []models.Habit
	if err := database.DB.Db.Where("completed = ?", false).Find(&incompleteHabits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(incompleteHabits)
}

// GetHabitStatistics calculates and returns statistics on habits
func GetHabitStatistics(c *fiber.Ctx) error {
	var totalHabitsCount int64
	var completedHabitsCount int64
	var incompleteHabitsCount int64

	// Get total number of habits
	if err := database.DB.Db.Model(&models.Habit{}).Count(&totalHabitsCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Get number of completed habits
	if err := database.DB.Db.Model(&models.Habit{}).Where("completed = ?", true).Count(&completedHabitsCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Calculate number of incomplete habits
	incompleteHabitsCount = totalHabitsCount - completedHabitsCount

	// Calculate completion rate
	var completionRate float64
	if totalHabitsCount > 0 {
		completionRate = float64(completedHabitsCount) / float64(totalHabitsCount) * 100
	}

	statistics := fiber.Map{
		"total_habits":            totalHabitsCount,
		"completed_habits":        completedHabitsCount,
		"incomplete_habits":       incompleteHabitsCount,
		"completion_rate_percent": completionRate,
	}

	return c.Status(200).JSON(statistics)
}
