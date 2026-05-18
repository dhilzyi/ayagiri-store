package api

import (
	"restaurant/internal/database"
)

type Handler struct {
	db       *database.Queries
	platform string
}

func NewHandler(db *database.Queries, platform string) *Handler {
	return &Handler{db: db, platform: platform}
}
