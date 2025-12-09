package models

import "gorm.io/gorm"

// Entity interface defines the contract for paginatable models
// Any model that implements this interface can use the generic Paginate function
type Entity interface {
	// Count returns the total number of records in the database for this entity
	Count(db *gorm.DB) int64

	// Take retrieves a paginated subset of records from the database
	// limit: maximum number of records to return
	// offset: number of records to skip for pagination
	Take(db *gorm.DB, limit int, offset int) interface{}
}
