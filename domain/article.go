package domain

import (
	"time"
)

// Article represents an article entity in the domain layer.
// This struct defines the core data structure and business rules for articles.
// IMPORTANT: This is part of the Domain layer - it contains ONLY:
//   - Data fields (no HTTP, no database queries, no frameworks)
//   - Business rules will be methods on this struct (e.g., Validate())
//
// This ensures the business logic is independent of how data is stored or presented.
type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    Author    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
