package controllers

import (
	"golangProject/database"
	"golangProject/models"
	"golangProject/util"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Register(c fiber.Ctx) error {
	// Create map to store request data
	var data map[string]string

	// Parse request body into data map
	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	// Check if password and password confirmation match
	if data["password"] != data["password_confirm"] {
		c.Status(400) // Set HTTP status to 400 Bad Request
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Passwords do not match",
		})
	}

	// Create new User instance with data from request
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    3,
	}

	// Generate a hashed password
	user.SetPassword(data["password"])

	// Save user to database
	database.DB.Create(&user)

	// Return the created user as JSON response
	return c.JSON(user)
}

func Login(c fiber.Ctx) error {
	// Create map to store request data
	var data map[string]string

	// Parse request body into data map
	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	// Create user variable to store query result
	var user models.User

	// Find user by email in database
	database.DB.Where("email = ?", data["email"]).First(&user)

	// Check if user was found (ID 0 means not found)
	if user.Id == 0 {
		c.Status(404) // Set HTTP status to 404 Not Found
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "email not found",
		})
	}

	// Compare provided password with stored hashed password using model method
	// This uses bcrypt.CompareHashAndPassword internally for secure comparison
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400) // Set HTTP status to 400 Bad Request
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "incorrect password",
		})
	}

	// Generate JWT token using the utility function - convert user ID to string
	token, err := util.GenerateJWT(strconv.Itoa(int(user.Id)))

	// Check if token creation failed
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Create HTTP-only cookie to store JWT
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	// Set the cookie in response
	c.Cookie(&cookie)

	// Return success message
	return c.JSON(fiber.Map{
		"message": "success login",
	})
}

func User(c fiber.Ctx) error {
	// Get JWT token from cookie
	cookie := c.Cookies("jwt")

	// Parse and validate JWT token, extract user ID
	id, _ := util.ParseJWT(cookie)

	// Create user variable to store query result
	var user models.User

	// Find user by ID from JWT claims
	database.DB.Where("id = ?", id).First(&user)

	// Return user information as JSON (excluding password due to json:"-" tag)
	return c.JSON(user)
}

func Logout(c fiber.Ctx) error {
	// Create a cookie with the same name as the JWT cookie but with empty value (remove cookie)
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	// Set the cookie in the response - this will overwrite and clear the existing JWT cookie
	c.Cookie(&cookie)

	// Return success message confirming logout
	return c.JSON(fiber.Map{
		"message": "success logout",
	})
}

// UpdateInfo allows authenticated users to update their personal information
// This includes first name, last name, and email address
// Users can only update their own information, identified via JWT token
func UpdateInfo(c fiber.Ctx) error {
	// Create a map to store the updated user data from request body
	var data map[string]string

	// Parse the JSON request body into the data map
	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	// Extract JWT token from the authentication cookie
	cookie := c.Cookies("jwt")

	// Parse the JWT token to get the user ID (issuer claim)
	id, _ := util.ParseJWT(cookie)

	// Convert the user ID string to integer
	userId, _ := strconv.Atoi(id)

	// Create a User instance with the authenticated user's ID and updated fields
	user := models.User{
		Id:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	// Update the user record in the database
	// This executes: UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;
	database.DB.Model(&user).Updates(user)

	// Return the updated user as JSON response
	return c.JSON(user)
}

// UpdatePassword allows authenticated users to change their password
// Requires password confirmation to prevent typos
// Users can only change their own password, identified via JWT token
func UpdatePassword(c fiber.Ctx) error {
	// Create a map to store password data from request body
	var data map[string]string

	// Parse the JSON request body into the data map
	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	// Validate that password and password confirmation match
	if data["password"] != data["password_confirm"] {
		c.Status(400) // Set HTTP status to 400 Bad Request
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Passwords do not match",
		})
	}

	// Extract JWT token from the authentication cookie
	cookie := c.Cookies("jwt")

	// Parse the JWT token to get the user ID (issuer claim)
	id, _ := util.ParseJWT(cookie)

	// Convert the user ID string to integer
	userId, _ := strconv.Atoi(id)

	// Create a User instance with only the ID for targeting the update
	user := models.User{
		Id: uint(userId),
	}

	// Hash the new password using the User model's SetPassword method
	user.SetPassword(data["password"])

	// Update only the password field in the database
	// This executes: UPDATE users SET password=? WHERE id=?;
	database.DB.Model(&user).Updates(user)

	// Return the updated user as JSON response (password excluded due to json:"-")
	return c.JSON(user)
}
