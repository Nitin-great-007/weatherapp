package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
)

type Weather struct {
	ID          int     `json:"id"`
	Location    string  `json:"location"`
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windspeed"`
	Condition   string  `json:"condition"`
}

func main() {
	// Open a connection to the MySQL database
	time.Sleep(9 * time.Second)
	db, err := sql.Open("mysql", "yourusername:yourpassword@tcp(mysql:3306)/yourdatabase")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new Fiber app
	app := fiber.New()

	// Endpoint to retrieve weather data by ID
	app.Get("/weather/:id", func(c *fiber.Ctx) error {
		weatherIDStr := c.Params("id")
		weatherID, err := strconv.Atoi(weatherIDStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid weather ID")
		}

		// Log the received ID for debugging
		fmt.Printf("Received weather ID: %d\n", weatherID)

		// Perform a query to retrieve the weather data from the database by ID
		row := db.QueryRow("SELECT id, location, date, temperature, windspeed, activities FROM weather_data WHERE id = ?", weatherID)
		var weather Weather
		if err := row.Scan(&weather.ID, &weather.Location, &weather.Date, &weather.Temperature, &weather.WindSpeed, &weather.Condition); err != nil {
			if err == sql.ErrNoRows {
				return c.Status(http.StatusNotFound).SendString("Weather data not found")
			}
			return c.Status(http.StatusInternalServerError).SendString("Internal server error")
		}

		// Convert the weather data to JSON and send it as the response
		return c.JSON(weather)
	})

	// Start the Fiber app
	port := 8081
	fmt.Printf("Server is running on port %d...\n", port)
	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}
