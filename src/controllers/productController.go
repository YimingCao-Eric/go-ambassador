package controllers

import (
	"golangProject/database"
	"golangProject/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// AllProducts retrieves a paginated list of products from the database
// This uses the generic Paginate function for consistent pagination
func AllProducts(c fiber.Ctx) error {
	// Extract the page number from query parameters, default to page 1 if not provided
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Use the generic Paginate function with Product entity
	// This provides standardized pagination response format
	return c.JSON(models.Paginate(database.DB, &models.Product{}, page))
}

// CreateProduct creates a new product in the database
// This function allows adding new products to the catalog
func CreateProduct(c fiber.Ctx) error {
	// Create a Product struct to hold the request data
	var product models.Product

	// Parse the JSON request body into the product struct
	if err := c.Bind().Body(&product); err != nil {
		return err
	}

	// Create the new product record in the database
	// This executes: INSERT INTO products (title, description, image, price) VALUES (?, ?, ?, ?);
	database.DB.Create(&product)

	// Return the created product as JSON response
	return c.JSON(product)
}

// GetProduct retrieves a specific product by ID from the database
// This is used for viewing individual product details
// URL: GET /api/products/:id
func GetProduct(c fiber.Ctx) error {
	// Extract the product ID from the URL parameter and convert to integer
	id, _ := strconv.Atoi(c.Params("id"))

	// Create a Product instance with the ID set for the database query
	product := models.Product{
		Id: uint(id),
	}

	// Find the product in the database by primary key (ID)
	// This executes: SELECT * FROM products WHERE id = ?;
	database.DB.Find(&product)

	// Return the product as JSON response
	return c.JSON(product)
}

// UpdateProduct updates an existing product's information
// This allows modifying product details like title, description, image, and price
// URL: PUT /api/products/:id
func UpdateProduct(c fiber.Ctx) error {
	// Extract the product ID from the URL parameter and convert to integer
	id, _ := strconv.Atoi(c.Params("id"))

	// Create a Product instance with the target ID
	product := models.Product{
		Id: uint(id),
	}

	// Parse the JSON request body into the product struct
	// This binds the updated fields to the product object
	if err := c.Bind().Body(&product); err != nil {
		return err
	}

	// Update the product record in the database
	// This executes: UPDATE products SET title=?, description=?, image=?, price=? WHERE id=?;
	database.DB.Model(&product).Updates(product)

	// Return the updated product as JSON response
	return c.JSON(product)
}

// DeleteProduct removes a product from the database
// This is a destructive operation and should be protected with proper authorization
// URL: DELETE /api/products/:id
func DeleteProduct(c fiber.Ctx) error {
	// Extract the product ID from the URL parameter and convert to integer
	id, _ := strconv.Atoi(c.Params("id"))

	// Create a Product instance with the target ID for deletion
	product := models.Product{
		Id: uint(id),
	}

	// Delete the product record from the database
	// This executes: DELETE FROM products WHERE id=?;
	database.DB.Delete(&product)

	// Return success status (204 No Content)
	return nil
}
