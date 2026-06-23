package main

import (
	"context"
	"os"
	"restaurant/internal/auth"
	"restaurant/internal/database"
)

func bootstrapAdmin(db *database.Queries) error {
	count, err := db.CountUsers(context.Background())
	if err != nil || count > 0 {
		return err
	}

	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return nil
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.CreateUser(context.Background(), database.CreateUserParams{
		Username:     username,
		PasswordHash: hash,
		Role:         "admin",
	})

	return err
}
