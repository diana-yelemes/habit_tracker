// handlers/facts.go
package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/diana-yelemes/habit_tracker/database"
	"github.com/diana-yelemes/habit_tracker/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllUserHabits retrieves all user habits
func GetAllUserHabits(c *fiber.Ctx) error {
	habits := []models.Habit{}
	database.DB.Db.Find(&habits)

	return c.Render("index", fiber.Map{
		"Title":    "Habit Tracker",
		"Subtitle": "Keep track on your daily habits!",
		"Habits":   habits,
	})
}

func NewHabitView(c *fiber.Ctx) error {
	return c.Render("new", fiber.Map{
		"Title":    "New Habit",
		"Subtitle": "Add a new daily habit to track!",
	})
}

type CreateHabitRequest struct {
	Name              string `json:"name" validate:"required"`
	RepeatCount       int    `json:"repeat_count"`
	TargetRepeatCount int    `json:"target_repeat_count"`
}

func CreateHabit(c *fiber.Ctx) error {
	var createHabitRequest CreateHabitRequest

	if err := c.BodyParser(&createHabitRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Create a new habit object
	habit := models.Habit{
		Name:              createHabitRequest.Name,
		RepeatCount:       0,
		TargetRepeatCount: createHabitRequest.TargetRepeatCount,
		Monday:            false,
		Tuesday:           false,
		Wednesday:         false,
		Thursday:          false,
		Friday:            false,
		Saturday:          false,
		Sunday:            false,
		Completed:         false,
		CompletionDate:    nil,
	}

	// Save the habit to the database
	if err := database.DB.Db.Create(&habit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating habit",
		})
	}

	return ConfirmationView(c)
}

func ConfirmationView(c *fiber.Ctx) error {
	return c.Render("confirmation", fiber.Map{
		"Title":    "Habit added successfully",
		"Subtitle": "Add more habits to the list!",
	})
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
	existingHabit.Name = updatedHabit.Name
	existingHabit.TargetRepeatCount = updatedHabit.TargetRepeatCount
	existingHabit.RepeatCount = updatedHabit.RepeatCount

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
func fetchAllHabits() ([]models.Habit, error) {
	var habits []models.Habit
	if err := database.DB.Db.Find(&habits).Error; err != nil {
		return nil, err
	}
	return habits, nil
}
func DeleteHabitView(c *fiber.Ctx) error {
	// Fetch the list of habits from the database
	habits, err := fetchAllHabits()

	if err != nil {
		// Handle the error (e.g., log, return an error response)
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching habits")
	}

	// Pass the habits to the template
	return c.Render("delete", fiber.Map{
		"Habits": habits,
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

	return c.JSON(fiber.Map{
		"message": "Habit deleted successfully",
	})
}

// GetHabitByID retrieves a single habit by ID
func GetHabitByID(c *fiber.Ctx) error {
	habitID := c.Params("id")
	var habit models.Habit

	// Include CalendarDays in the query
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
		sql += fmt.Sprintf(" AND name ILIKE '%%%s%%'", name)
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

type UpdateRepeatCountRequest struct {
	DayIndex  int  `json:"day_index"`
	Completed bool `json:"completed"`
}

type UpdateRepeatCountResponse struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	RepeatCount       int        `json:"repeat_count"`
	TargetRepeatCount int        `json:"target_repeat_count"`
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

// UpdateRepeatCount updates the repeat count and completed status of a habit based on the day index
func UpdateRepeatCount(c *fiber.Ctx) error {
	// Parse habit ID from the URL parameters
	habitID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid habit ID",
		})
	}

	// Parse request body
	var requestBody UpdateRepeatCountRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Retrieve habit from the database
	var habit models.Habit
	result := database.DB.Db.First(&habit, habitID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Habit not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving habit",
		})
	}

	// Update repeat count based on day index
	switch requestBody.DayIndex {
	case 0:
		habit.Monday = requestBody.Completed
	case 1:
		habit.Tuesday = requestBody.Completed
	case 2:
		habit.Wednesday = requestBody.Completed
	case 3:
		habit.Thursday = requestBody.Completed
	case 4:
		habit.Friday = requestBody.Completed
	case 5:
		habit.Saturday = requestBody.Completed
	case 6:
		habit.Sunday = requestBody.Completed
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid day index",
		})
	}

	// Calculate the new repeat count based on the updated days of the week
	repeatCount := 0
	if habit.Monday {
		repeatCount++
	}
	if habit.Tuesday {
		repeatCount++
	}
	if habit.Wednesday {
		repeatCount++
	}
	if habit.Thursday {
		repeatCount++
	}
	if habit.Friday {
		repeatCount++
	}
	if habit.Saturday {
		repeatCount++
	}
	if habit.Sunday {
		repeatCount++
	}

	// Update the repeat count
	habit.RepeatCount = repeatCount

	// Update the completed status based on the 'completed' parameter
	habit.Completed = requestBody.Completed

	// Save the updated habit to the database
	if err := database.DB.Db.Save(&habit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating habit",
		})
	}

	response := UpdateRepeatCountResponse{
		ID:                habit.ID,
		Name:              habit.Name,
		RepeatCount:       habit.RepeatCount,
		TargetRepeatCount: habit.TargetRepeatCount,
		Monday:            habit.Monday,
		Tuesday:           habit.Tuesday,
		Wednesday:         habit.Wednesday,
		Thursday:          habit.Thursday,
		Friday:            habit.Friday,
		Saturday:          habit.Saturday,
		Sunday:            habit.Sunday,
		Completed:         habit.Completed,
		CompletionDate:    habit.CompletionDate,
	}

	return c.JSON(response)
}
