package controllers

import (
	"golangProject/database"
	"golangProject/middlewares"
	"golangProject/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// AllUsers retrieves a paginated list of users from the database
// This uses the generic Paginate function for consistent pagination
// Security Note: In production, consider adding role-based access control
func AllUsers(c fiber.Ctx) error {
	if err := middlewares.IsAuthorized(c, "users"); err != nil {
		return err
	}

	// Extract the page number from query parameters, default to page 1 if not provided
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Use the generic Paginate function with User entity
	// This provides standardized pagination response format
	return c.JSON(models.Paginate(database.DB, &models.User{}, page))
}

// CreateUser creates a new user with a default password
// This function allows creating users programmatically (admin function)
// In production, consider adding validation and proper error handling
func CreateUser(c fiber.Ctx) error {
	if err := middlewares.IsAuthorized(c, "users"); err != nil {
		return err
	}

	// Create a User struct to hold the request data
	var user models.User

	// Parse the JSON request body into the user struct
	if err := c.Bind().Body(&user); err != nil {
		return err
	}

	// Set a default password using the User model's SetPassword method
	// This encapsulates the password hashing logic within the model
	// Note: "3" is a hardcoded default password - consider making this configurable
	user.SetPassword("3")

	// Create the new user record in the database
	// This executes: INSERT INTO users (...) VALUES (...);
	database.DB.Create(&user)

	// Return the created user as JSON response (password excluded)
	return c.JSON(user)
}

// GetUser retrieves a specific user by ID from the database
// This is typically used for viewing individual user profiles
// URL: GET /api/users/:id
func GetUser(c fiber.Ctx) error {
	if err := middlewares.IsAuthorized(c, "users"); err != nil {
		return err
	}

	// Extract the user ID from the URL parameter and convert to integer
	// c.Params("id") gets the ":id" value from the route
	id, _ := strconv.Atoi(c.Params("id"))

	// Create a User instance with the ID set for the database query
	user := models.User{
		Id: uint(id),
	}

	// Find the user in the database by primary key (ID)
	// This executes: SELECT * FROM users WHERE id = ?;
	database.DB.Preload("Role").Find(&user)

	// Return the user as JSON response (password excluded due to json:"-")
	return c.JSON(user)
}

// UpdateUser updates an existing user's information
// This allows modifying user details like name and email
// URL: PUT /api/users/:id
func UpdateUser(c fiber.Ctx) error {
	if err := middlewares.IsAuthorized(c, "users"); err != nil {
		return err
	}

	// Extract the user ID from the URL parameter and convert to integer
	id, _ := strconv.Atoi(c.Params("id"))

	// Create a User instance with the target ID
	user := models.User{
		Id: uint(id),
	}

	// Parse the JSON request body into the user struct
	// This binds the updated fields (first_name, last_name, email) to the user object
	if err := c.Bind().Body(&user); err != nil {
		return err
	}

	// Update the user record in the database
	// .Model() specifies which record to update, .Updates() applies the changes
	// This executes: UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;
	database.DB.Model(&user).Updates(user)

	// Return the updated user as JSON response
	return c.JSON(user)
}

// DeleteUser removes a user from the database
// This is a destructive operation and should be protected with proper authorization
// URL: DELETE /api/users/:id
func DeleteUser(c fiber.Ctx) error {
	if err := middlewares.IsAuthorized(c, "users"); err != nil {
		return err
	}

	// Extract the user ID from the URL parameter and convert to integer
	id, _ := strconv.Atoi(c.Params("id"))

	// Create a User instance with the target ID for deletion
	user := models.User{
		Id: uint(id),
	}

	// Delete the user record from the database
	// This executes: DELETE FROM users WHERE id=?;
	database.DB.Delete(&user)

	// Return nil (no content) to indicate successful deletion
	// Alternatively, you could return a success message
	return nil
}
