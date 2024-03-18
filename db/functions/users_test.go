package functions

import (
	"context"
	"errors"
	"log"
	"shopifyx/configs"
	"shopifyx/db/connections"
	"shopifyx/db/entity"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MockPool implements pgxpool.Pool interface for testing purposes.
type MockPool struct{}

func (mp *MockPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return nil, nil // Implement this if needed for your tests
}

func (mp *MockPool) Close() {}

func TestRegister(t *testing.T) {
	// Mocking dependencies
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	mockPool, err := connections.NewPgConn(config)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer mockPool.Close()

	mockConfig := configs.Config{} // You might need to fill this with necessary fields for testing

	// Creating user instance
	user := NewUser(mockPool, mockConfig)

	// Test case 1: Registering a new user successfully
	usr := entity.User{
		Name:            "John Doe",
		CredentialType:  "email",
		CredentialValue: "john@example.com",
		Password:        "password123",
	}
	createdUser, err := user.Register(context.Background(), usr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdUser.Id == "" {
		t.Error("Expected user ID to be set, got an empty string")
	}

	// Test case 2: Registering with an existing username
	existingUser := entity.User{
		Name:            "Jane Doe",
		CredentialType:  "email",
		CredentialValue: "jane@example.com",
		Password:        "password123",
	}
	_, err = user.Register(context.Background(), existingUser)
	if !errors.Is(err, errors.New("EXISTING_USERNAME")) {
		t.Errorf("Expected error 'EXISTING_USERNAME', got %v", err)
	}

	// Additional test cases can be added as needed
}

func TestLogin(t *testing.T) {
	// Mocking dependencies
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	mockPool, err := connections.NewPgConn(config)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer mockPool.Close()

	mockConfig := configs.Config{} // You might need to fill this with necessary fields for testing

	// Creating user instance
	user := NewUser(mockPool, mockConfig)

	// Test case 1: Logging in with existing user and correct password
	existingUser := entity.User{
		CredentialType:  "email",
		CredentialValue: "john@example.com",
		Password:        "password123",
	}
	result, err := user.Login(context.Background(), existingUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Id == "" {
		t.Error("Expected user ID to be set, got an empty string")
	}

	// Test case 2: Logging in with non-existing user
	nonExistingUser := entity.User{
		CredentialType:  "email",
		CredentialValue: "nonexisting@example.com",
		Password:        "password123",
	}
	_, err = user.Login(context.Background(), nonExistingUser)
	if !errors.Is(err, errors.New("USER_NOT_FOUND")) {
		t.Errorf("Expected error 'USER_NOT_FOUND', got %v", err)
	}

	// Test case 3: Logging in with existing user and incorrect password
	incorrectPasswordUser := entity.User{
		CredentialType:  "email",
		CredentialValue: "john@example.com",
		Password:        "incorrectpassword",
	}
	_, err = user.Login(context.Background(), incorrectPasswordUser)
	if !errors.Is(err, errors.New("INVALID_PASSWORD")) {
		t.Errorf("Expected error 'INVALID_PASSWORD', got %v", err)
	}

	// Additional test cases can be added as needed
}
