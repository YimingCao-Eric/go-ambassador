package models

import (
	"math"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Paginate provides a generic pagination solution for any Entity
// This function eliminates code duplication for paginated endpoints
// Parameters:
//   - db: GORM database connection
//   - entity: Entity interface implementation (User, Product, etc.)
//   - page: current page number for pagination
//
// Returns:
//   - fiber.Map: standardized pagination response with data and metadata
func Paginate(db *gorm.DB, entity Entity, page int) fiber.Map {
	// Set the number of records to display per page
	limit := 5

	// Calculate the offset for database query
	offset := (page - 1) * limit

	// Retrieve paginated data using the entity's Take method
	data := entity.Take(db, limit, offset)

	// Get total record count using the entity's Count method
	total := entity.Count(db)

	// Return standardized pagination response
	return fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(total) / float64(limit)), // Calculate the last page number for pagination navigation
		},
	}
}
