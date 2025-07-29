package main

import (
	"fmt"
	"log"

	"github.com/harungurubudi/mtsg/internal/di"
)

func main() {
	// Initialize the authentication use case using Wire
	authUseCase := di.InitializeAuthUseCase()

	// Test the DI setup
	fmt.Println("✅ Dependency injection setup successful!")
	fmt.Printf("Auth use case type: %T\n", authUseCase)

	// You can now use authUseCase for authentication operations
	// Example: authUseCase.Login(ctx, credential)

	log.Println("MTSG application started with Wire DI")
}
