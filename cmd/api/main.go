// Vistara Backend API Server
// This is the main entry point for the Vistara Backend application.
// It initializes all necessary components including database connections,
// HTTP server, middleware, and routing configuration.
package main

import "github.com/vistara-studio/vistara-be/internal/bootstrap"

func main() {
	// Initialize the application with all dependencies
	// This includes database setup, HTTP server configuration,
	// middleware registration, and route mounting
	if err := bootstrap.Initialize(); err != nil {
		panic(err)
	}
}
