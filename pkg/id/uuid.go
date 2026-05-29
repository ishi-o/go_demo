package id

import (
	"github.com/google/uuid"
)

// Generate creates a new UUID v7
func Generate() string {
	return uuid.New().String()
}

// MustGenerate creates a new UUID and panics on error
func MustGenerate() string {
	return uuid.New().String()
}

// Validate checks if the given string is a valid UUID
func Validate(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
