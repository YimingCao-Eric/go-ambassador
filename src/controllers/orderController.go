package controllers

import (
	"encoding/csv"
	"golangProject/database"
	"golangProject/models"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// AllOrders returns a paginated list of all orders with their items
// Uses the generic Paginate function for consistent pagination
// Query parameter: page (defaults to 1) - specifies which page of results to return
func AllOrders(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	return c.JSON(models.Paginate(database.DB, &models.Order{}, page))
}

// Export generates a CSV file containing all orders and order items
// Creates a structured export suitable for spreadsheets or data analysis
// The CSV file is saved temporarily and sent as a download to the client
func Export(c fiber.Ctx) error {
	filePath := "./csv/order.csv"

	// Create the CSV file and populate with order data
	if err := CreateFile(filePath); err != nil {
		return err
	}

	// Send the generated file to client as a download
	return c.Download(filePath)
}

// CreateFile generates the actual CSV file with order data
// File structure:
//   - Header row with column names
//   - Each order spans multiple rows: one main row + one per order item
//   - Empty cells used for visual grouping of order items under their order
func CreateFile(filePath string) error {
	// Create the CSV file
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer file.Close() // Ensure file is closed when function exits

	// Create CSV writer with buffering
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written to file

	var orders []models.Order

	// Load all orders with their associated items from database
	database.DB.Preload("OrderItems").Find(&orders)

	// Write CSV header row
	writer.Write([]string{
		"ID", "Name", "Email", "Product Title", "Price", "Quantity",
	})

	// Write order data to CSV
	for _, order := range orders {
		// Write the main order row (shows order info, empty product columns)
		data := []string{
			strconv.Itoa(int(order.Id)),
			order.FirstName + " " + order.LastName,
			order.Email,
			"",
			"",
			"",
		}
		if err := writer.Write(data); err != nil {
			return err
		}

		// Write each order item as a separate row
		for _, orderItem := range order.OrderItems {
			data := []string{
				"",
				"",
				"",
				orderItem.ProductTitle,
				strconv.Itoa(int(orderItem.Price)),
				strconv.Itoa(int(orderItem.Quantity)),
			}
			if err := writer.Write(data); err != nil {
				return err
			}
		}
	}
	return nil
}

// Sales represents daily sales data for chart visualization
// Used by the Chart endpoint to return sales trends over time
type Sales struct {
	Date string `json:"date"`
	Sum  string `json:"sum"`
}

// Chart returns daily sales data for visualization
// Uses raw SQL to group sales by date and calculate daily totals
// Returns data suitable for line charts or sales trend analysis
func Chart(c fiber.Ctx) error {
	var sales []Sales

	// Execute raw SQL query to get daily sales totals
	// Groups orders by creation date and sums the product of price * quantity
	database.DB.Raw(`
		SELECT DATE_FORMAT(o.create_at, '%Y-%m-%d') as date, SUM(oi.price*oi.quantity) as sum 
		FROM orders o 
		JOIN order_items oi on o.id=oi.order_id 
		GROUP BY date
		`).Scan(&sales)
	return c.JSON(sales)
}
