// @title Forum API
// @version 1.0
// @description API for Forum-Go project
// @termsOfService http://example.com/terms/

// @contact.name Tolik
// @contact.url http://example.com
// @contact.email rediska1203@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"fmt"
	"forum1/internal/app"
)

func main() {
	fmt.Println("Starting forum application...")
	app.Run()
}
